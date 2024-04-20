package executors

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"multivac.network/services/agents/graph/model"
	"multivac.network/services/agents/messages"
	"multivac.network/services/agents/providers"
	"net/http"
	"net/url"
	"strings"
	"text/template"
)

import "embed"

//go:embed embedded/prompts/*
var embeddings embed.FS

type Agent struct {
	Descriptor        *model.Agent
	ReplyChannel      chan *messages.WebSocketMessage
	CommandChannel    chan<- *messages.CommandType
	Input             chan *providers.Message
	Internal          chan *providers.Message
	prompt            string
	thoughtPrompt     string
	defaultPrompt     string
	functionPrompt    string
	Thought           string
	EvaluationChannel chan *providers.Message
	Context           []providers.Message
	ThoughtContext    []providers.Message
	service           providers.ModelProvider
}

func NewAgent(service providers.ModelProvider, agent *model.Agent, output chan *messages.WebSocketMessage, input chan *providers.Message) *Agent {
	result := &Agent{
		Descriptor:     agent,
		Internal:       make(chan *providers.Message),
		Input:          input,
		prompt:         agent.Prompt,
		service:        service,
		Context:        make([]providers.Message, 0),
		ReplyChannel:   output,
		CommandChannel: make(chan<- *messages.CommandType),
	}
	go result.initialize()
	go result.initializeEvaluation()
	defaultPrompt, err := embeddings.ReadFile("embedded/prompts/default")

	if err != nil {
		panic(err)
	}
	log.Println(fmt.Sprintf("Agent Prompt: %s", agent.Prompt))
	result.Context = append(result.Context, providers.Message{Role: "system", Content: agent.Prompt})
	result.defaultPrompt = string(defaultPrompt)
	return result
}

func (agent *Agent) Chat(context string, text string) (err error) {

	agent.Internal <- &providers.Message{Role: "user", Content: text}
	return nil
}

func fetchInformation(index string, text string) string {
	client := &http.Client{}
	query := url.QueryEscape(text)
	response, err := client.Get("http://multivac-embeddings-service.default.svc.cluster.local/context?index=" + index + "&q=" + query)
	if err != nil {
		log.Println(err)
		return ""
	}
	content, err := io.ReadAll(response.Body)
	return string(content)
}

func (agent *Agent) responseHandler(message providers.Message) {
	agent.Context = append(agent.Context, message)
	agent.ReplyChannel <- messages.Message("chat-response", messages.ReplyMessage{
		Agent:   agent.Descriptor.Name,
		Content: message.Content,
	})
	request := providers.Request{Messages: make([]providers.Message, 0), Stream: false}

	content, err := generateEvaluationTemplate(struct {
		Group       string
		Description string
		Message     string
		Agent       string
	}{
		Agent:   agent.Descriptor.Name,
		Message: message.Content,
	})
	request.Messages = append(request.Messages, providers.Message{
		Role:    "user",
		Content: content,
	})

	err = agent.service.SendRequest(request, agent.EvaluationChannel)
	if err != nil {
		log.Println(err)
	}
}

func generateEvaluationTemplate(data interface{}) (string, error) {
	t, err := template.New("group-template").Parse(`
		You have received a response from an agent. The agent is '{{.Agent}}'. The agent responded with '{{.Message}}'.
		You only respond with either the word 'true' or 'false'. If the agent's response
		is correct and no agent in the list below is needed to add to the answer, respond with 'true'. 
		If the agent's response is incorrect or requires another agent to also respond, respond with 'false'.
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

func (agent *Agent) initialize() {
	for {
		select {
		case message := <-agent.Internal:
			agent.responseHandler(*message)
		}
	}
}

func (agent *Agent) initializeEvaluation() {
	for {
		select {
		case message := <-agent.EvaluationChannel:
			complete := strings.Contains(message.Content, "true")
			if !complete {
				agent.Input <- &providers.Message{Role: "user", Content: message.Content}
			}
		}
	}
}
