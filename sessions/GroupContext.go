package sessions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"multivac.network/services/agents/data"
	"multivac.network/services/agents/executors"
	"multivac.network/services/agents/graph/model"
	"multivac.network/services/agents/messages"
	"multivac.network/services/agents/providers"
	"strings"
	"text/template"
)

type GroupContext struct {
	group              *data.GroupModel
	agents             []*executors.Agent
	socket             *websocket.Conn
	model              providers.ModelProvider
	routing            chan *providers.Message
	evaluationResponse chan *providers.Message
	Evaluation         chan *providers.Message
	Input              chan *providers.Message
	Output             chan *messages.WebSocketMessage
	Status             chan *messages.WebSocketMessage
	messages           []providers.Message
}

func NewGroupContext(group *data.GroupModel, socket *websocket.Conn, provider providers.ModelProvider, agents []*model.Agent) *GroupContext {
	agentExecutors := make([]*executors.Agent, 0)
	outputChannel := make(chan *messages.WebSocketMessage)
	inputChannel := make(chan *providers.Message)
	evaluationChannel := make(chan *providers.Message)
	evaluationResponseChannel := make(chan *providers.Message)
	for _, agent := range agents {
		agentExecutors = append(agentExecutors, executors.NewAgent(provider, agent, outputChannel, evaluationChannel))
	}
	result := &GroupContext{
		group:              group,
		agents:             agentExecutors,
		socket:             socket,
		model:              provider,
		messages:           make([]providers.Message, 0),
		routing:            make(chan *providers.Message),
		Input:              inputChannel,
		Output:             outputChannel,
		Evaluation:         evaluationChannel,
		evaluationResponse: evaluationResponseChannel,
	}
	go result.start()
	return result
}

func (group *GroupContext) start() {
	go group.initializeRouting()
	go group.initializeOutput()
	go group.initializeInput()
	go group.initializeSocket()
}

func (c *GroupContext) initializeInput() {
	agentDescriptions := strings.Builder{}
	for _, agent := range c.agents {
		agentDescriptions.WriteString(fmt.Sprintf("%s:%s %s\n", agent.Descriptor.ID, agent.Descriptor.Name, agent.Descriptor.Description))
	}

	for {
		select {
		case message := <-c.Input:
			request := providers.Request{Messages: make([]providers.Message, 0), Stream: false}

			content, err := generateTemplate(struct {
				Group       string
				Description string
				Message     string
				Agents      string
			}{
				Group:       c.group.Name,
				Description: c.group.Description,
				Message:     message.Content,
				Agents:      agentDescriptions.String(),
			})

			request.Messages = append(request.Messages, providers.Message{
				Role:    "user",
				Content: content,
			})

			err = c.model.SendRequest(request, c.routing)
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}
}

func generateTemplate(data interface{}) (string, error) {
	t, err := template.New("group-template").Parse(`
		You are a router for a group called '{{.Group}}'. The group is a collection of agents that are 
		working together to solve a problem. The group is described as '{{.Description}}'. The message 
		"{{.Message}}" was received by the group. The agents in the group are '{{.Agents}}'. 
		Decide which agents should respond and to what prompt with a score between 0 and 1 of how confident you are they 
		are the right agent. Confidence scores should be based on the description of the agent relative to the request. 
		Higher scores are more relevant agents than lower.
		Respond with a JSON array of {"id": "<agent id>", "prompt": "<prompt>", "confidence": <confidence score>} pairs. 
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

func (c *GroupContext) initializeOutput() {
	for {
		select {
		case message := <-c.Output:
			output, err := json.Marshal(message)
			err = c.emitMessage(string(output))
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (c *GroupContext) emitMessage(message string) error {
	return c.socket.WriteMessage(websocket.TextMessage, []byte(message))
}

func (c *GroupContext) Close() error {
	return c.socket.Close()
}

func (c *GroupContext) initializeSocket() {
	for {
		_, p, err := c.socket.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		message := providers.Message{}
		err = json.Unmarshal(p, &message)
		if err != nil {
			log.Println(err)
			continue
		}
		c.Input <- &message
	}
}

func (c *GroupContext) initializeRouting() {
	for {
		select {
		case message := <-c.routing:
			var agents []AgentSelection
			err := json.Unmarshal([]byte(message.Content), &agents)
			if err != nil {
				log.Println(err)
				continue
			}
			for _, agent := range agents {
				if agent.Confidence < 0.8 {
					continue
				}
				for _, a := range c.agents {
					if a.Descriptor.ID != agent.Id {
						continue
					}
					if err != nil {
						log.Println(err)
						continue
					}

					c.Output <- messages.Message("status", StatusMessage{
						Status:  "typing",
						Content: fmt.Sprintf("%s is responding", a.Descriptor.Name)})

					a.Chat("", agent.Prompt)
				}
			}
			fmt.Printf("Internal message: %s\n", message.Content)
		}
	}
}

type StatusMessage struct {
	Status  string `json:"status"`
	Content string `json:"content"`
}

type AgentSelection struct {
	Id         string  `json:"id"`
	Prompt     string  `json:"prompt"`
	Confidence float64 `json:"confidence"`
}
