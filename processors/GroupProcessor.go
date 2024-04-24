package processors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"multivac.network/services/agents/data"
	"multivac.network/services/agents/messages"
	"multivac.network/services/agents/providers"
	"strings"
	"text/template"
)

type GroupProcessor struct {
	*Input[*messages.ConversationMessage]
	Memory      *data.VectorStore
	Loopback    *Input[*messages.AgentMessage]
	FinalOutput *Output[*messages.AgentMessage]
	Model       *data.GroupModel

	User         string
	Context      []*messages.ConversationMessage
	provider     providers.ModelProvider
	agents       []*AgentProcessor
	descriptions string
}

// CommitContext saves the context of the group to the vector database.

// NewGroupProcessor creates a new group processor
func NewGroupProcessor(user string, group *data.GroupModel, provider providers.ModelProvider) *GroupProcessor {
	processor := &GroupProcessor{
		Model:       group,
		Memory:      data.NewVectorStore(group.ID),
		User:        user,
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
	// processor.Memory.Clear()
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

	// add up to the last 4 messages to the context
	count := 0
	if len(gp.Context)-4 > 0 {
		count = len(gp.Context) - 4
	}
	for _, context := range gp.Context[count:] {
		request.Messages = append(request.Messages, providers.Message{
			Role:    context.Role,
			Content: context.Content,
		})
	}
	request.Messages = append(request.Messages, providers.Message{
		Role:    "system",
		Content: content,
	})
	request.Messages = append(request.Messages, providers.Message{
		Role:    "user",
		Content: message.Content,
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
	// templates should be moved to the database.
	t, err := template.New("group-template").Parse(`
		You are a conversation router for a group called '{{.Group}}'. 
		The group is a collection of agents that are working together to solve a problem. 
		The group is described as '{{.Description}}'. 
		The agents in the group are: 
			'{{.Agents}}'
		Decide which agents should respond and to what prompt with a score between 0 and 1 with 1 being the most confident you are they
		are the right agent to respond and 0 being the least. 
		RULES:
			1. If an agent is referenced by name. Then that agent, and ONLY that agent, should have a confidence score of 1.
			2. Confidence scores should be based on the description of the agent relative to the request unless mentioned by name.
			3. Scores lower than 0.8 should be excluded from your result.
			4. If any agent has a score of 1, then only that agent should respond.
			5. Respond with a JSON array in the following format: {"id": "<agent id>","name":"<agent name>", "prompt": "<prompt>", "confidence": <confidence score>}.
			6. Respond only with the proper formatted JSON. Copy the original prompt for each agent you want to respond.
			7. THE RESPONSE SHOULD BE JSON ONLY.
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

			message.Context = append(message.Context, &messages.ConversationMessage{Role: "user", Content: message.Content})
			message.Content = fmt.Sprintf("<Agent>%s</Agent> %s", gp.User, message.Content)
			a.input <- message
		}
	}
	_ = gp.Memory.Commit("<Agent>" + gp.User + "</Agent>" + message.Content)
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
				err := gp.Memory.Commit(conversationMessage.Content)
				if err != nil {
					log.Printf("error committing context: %s", err)
				}
				continue
			}
		}
	}()
}

type AgentSelection struct {
	Id         string  `json:"id"`
	Name       string  `json:"name"`
	Prompt     string  `json:"prompt"`
	Confidence float64 `json:"confidence"`
}
