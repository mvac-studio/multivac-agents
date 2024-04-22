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
	Loopback     *Input[*messages.AgentMessage]
	FinalOutput  *Output[*messages.AgentMessage]
	Model        *data.GroupModel
	Context      []*messages.ConversationMessage
	provider     providers.ModelProvider
	agents       []*AgentProcessor
	descriptions string
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
		are the right agent. Confidence scores should be based on the description of the agent relative to the request.
		Higher scores are more relevant agents than lower.
		Respond with a JSON array of {"id": "<agent id>","name":"<agent name>", "prompt": "<prompt>", "confidence": <confidence score>} pairs.
		Respond only with the proper formatted JSON.
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
	message.Context = gp.Context
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

			message.Context = append(message.Context, &messages.ConversationMessage{Role: "user", Content: agent.Prompt})
			a.input <- message
		}
	}
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
				gp.Context = append(gp.Context, &messages.ConversationMessage{
					Role:    "assistant",
					Content: fmt.Sprintf("<Agent>%s</Agent> %s", message.Agent, message.Content)})
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