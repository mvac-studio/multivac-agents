package processors

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	chromago "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/collection"
	"github.com/amikos-tech/chroma-go/types"
	"log"
	"multivac.network/services/agents/data"
	"multivac.network/services/agents/messages"
	"multivac.network/services/agents/providers"
	"strings"
	"text/template"
)

type GroupProcessor struct {
	*Input[*messages.ConversationMessage]
	vectorClient *chromago.Client
	Loopback     *Input[*messages.AgentMessage]
	FinalOutput  *Output[*messages.AgentMessage]
	Model        *data.GroupModel
	Context      []*messages.ConversationMessage
	provider     providers.ModelProvider
	agents       []*AgentProcessor
	descriptions string
}

// CommitContext saves the context of the group to the vector database.
func (gp *GroupProcessor) CommitContext(conversationMessage *messages.ConversationMessage) error {
	embedfn := types.NewConsistentHashEmbeddingFunction()
	col, err := gp.vectorClient.NewCollection(context.Background(),
		collection.WithName(gp.Model.ID),
		collection.WithEmbeddingFunction(embedfn),
		collection.WithCreateIfNotExist(true),
	)
	if err != nil {
		return err
	}

	rs, err := types.NewRecordSet(types.WithEmbeddingFunction(embedfn), types.WithIDGenerator(types.NewULIDGenerator()))
	rs.WithRecord(types.WithDocument(conversationMessage.Content))
	if err != nil {
		log.Println(err)
	}

	_, err = rs.BuildAndValidate(context.TODO())
	_, err = col.AddRecords(context.Background(), rs)
	if err != nil {
		return err
	}
	return nil
}

// NewGroupProcessor creates a new group processor
func NewGroupProcessor(group *data.GroupModel, provider providers.ModelProvider) *GroupProcessor {
	processor := &GroupProcessor{
		Model:       group,
		Context:     make([]*messages.ConversationMessage, 0),
		FinalOutput: NewOutputProcessor[*messages.AgentMessage](),
		provider:    provider,
		agents:      make([]*AgentProcessor, 0),
	}

	processor.Input = NewInputProcessor[*messages.ConversationMessage]()
	processor.Loopback = NewInputProcessor[*messages.AgentMessage]()
	processor.FinalOutput = NewOutputProcessor[*messages.AgentMessage]()
	processor.initialize()
	processor.initializeLoopback()
	client, err := chromago.NewClient("http://chromadb-service.default.svc.cluster.local:8000")
	if err != nil {
		log.Println(err)
	}
	processor.vectorClient = client

	// TODO: this is temporary while working out the issues.
	// processor.vectorClient.DeleteCollection(context.Background(), group.ID)
	return processor
}

// AddAgent adds an agent to the group
func (gp *GroupProcessor) AddAgent(agent *AgentProcessor) error {
	agent.To(gp.Loopback)
	gp.agents = append(gp.agents, agent)
	gp.updateDescriptions()
	return nil
}

// Process processes the message
func (gp *GroupProcessor) Process(message *messages.ConversationMessage) error {
	request := providers.Request{Messages: make([]providers.Message, 0), Stream: false}

	content, err := generateTemplate(struct {
		Group       string
		Description string
		Message     string
		Agents      string
	}{
		Group:       gp.Model.Name,
		Description: gp.Model.Description,
		Message:     message.Content,
		Agents:      gp.descriptions,
	})

	request.Messages = append(request.Messages, providers.Message{
		Role:    "user",
		Content: content,
	})

	response, err := gp.provider.SendRequest(request)
	if err != nil {
		log.Println(err)
		return err
	}
	gp.route(message, response)
	return err
}

func generateTemplate(data interface{}) (string, error) {
	t, err := template.New("group-template").Parse(`
		You are a router for a group called '{{.Group}}'. The group is a collection of agents that are
		working together to solve a problem. The group is described as '{{.Description}}'. The message
		"{{.Message}}" was received by the group. The agents in the group are '{{.Agents}}'.
		Decide which agents should respond and to what prompt with a score between 0 and 1 of how confident you are they
		are the right agent. If an agent is mentioned by name with '@name'. Then only that agent should respond.
		Confidence scores should be based on the description of the agent relative to the request.
		Higher scores are more relevant agents than lower.
		Respond with a JSON array of {"id": "<agent id>","name":"<agent name>", "prompt": "<prompt>", "confidence": <confidence score>} pairs.
		Respond only with the proper formatted JSON. Reword the prompt for each agent into question form expanding or editing the prompt
		to be detailed for the specific agent that is handling the request. Do not ask them about themselves.
	`)
	if err != nil {
		return "", err
	}

	buffer := &bytes.Buffer{}

	err = t.Execute(buffer, data)
	if err != nil {
		log.Println(err)
	}
	return buffer.String(), err
}

func (gp *GroupProcessor) updateDescriptions() {
	builder := strings.Builder{}
	for _, agent := range gp.agents {
		builder.WriteString(fmt.Sprintf("%s:%s %s\n", agent.AgentModel.ID, agent.AgentModel.Name, agent.AgentModel.Description))
	}
	gp.descriptions = builder.String()
}

func (gp *GroupProcessor) route(message *messages.ConversationMessage, response *providers.Message) {
	var agents []AgentSelection
	err := json.Unmarshal([]byte(response.Content), &agents)
	var ctx []*messages.ConversationMessage
	if len(gp.Context) > 6 {
		ctx = gp.Context[len(gp.Context)-2:]
	} else {
		ctx = gp.Context

	}

	message.Context = ctx
	reference := gp.getContext(message)
	for _, agent := range agents {
		if agent.Confidence < 0.8 {
			continue
		}
		for _, a := range gp.agents {
			if a.AgentModel.ID != agent.Id {
				continue
			}
			if err != nil {
				log.Println(err)
				continue
			}

			prompt := fmt.Sprintf("<Reference>%s</Reference> %s", reference, agent.Prompt)
			message.Context = append(message.Context, &messages.ConversationMessage{Role: "user", Content: prompt})
			a.input <- message
		}
	}
	_ = gp.CommitContext(message)
}

func (gp *GroupProcessor) initialize() {
	go func() {
		for {
			select {
			case message := <-gp.input:
				err := gp.Process(message)
				if err != nil {
					log.Println(err)
					continue
				}
			}
		}
	}()
}

// initializeLoopback initializes the loopback channel for processing
func (gp *GroupProcessor) initializeLoopback() {
	go func() {
		for {
			select {
			case message := <-gp.Loopback.input:
				gp.FinalOutput.output <- message
				conversationMessage := &messages.ConversationMessage{
					Role:    "assistant",
					Content: fmt.Sprintf("<Agent>%s</Agent> %s", message.Agent, message.Content)}

				gp.Context = append(gp.Context, conversationMessage)
				err := gp.CommitContext(conversationMessage)
				if err != nil {
					log.Printf("error committing context: %s", err)
				}
				continue
			}
		}
	}()
}

func (gp *GroupProcessor) getContext(message *messages.ConversationMessage) string {
	embedfn := types.NewConsistentHashEmbeddingFunction()
	col, err := gp.vectorClient.NewCollection(context.Background(),
		collection.WithName(gp.Model.ID),
		collection.WithEmbeddingFunction(embedfn),
		collection.WithCreateIfNotExist(true),
	)
	if err != nil {
		return ""
	}
	result, err := col.Query(context.TODO(), []string{message.Content}, 3, nil, nil, nil)
	if err != nil {
		log.Println(err)
	}
	builder := strings.Builder{}
	for _, r := range result.Documents {
		for _, d := range r {
			builder.WriteString(d)
		}
	}
	return builder.String()
}

type AgentSelection struct {
	Id         string  `json:"id"`
	Name       string  `json:"name"`
	Prompt     string  `json:"prompt"`
	Confidence float64 `json:"confidence"`
}
