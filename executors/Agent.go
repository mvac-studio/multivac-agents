package executors

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"multivac.network/services/agents/graph/model"
	"multivac.network/services/agents/messages"
	"multivac.network/services/agents/services"
	"net/http"
	"net/url"
	"text/template"
)

import "embed"

//go:embed embedded/prompts/*
var embeddings embed.FS

type Agent struct {
	description    *model.Agent
	ReplyChannel   chan *messages.WebSocketMessage
	CommandChannel chan<- *messages.CommandType
	prompt         string
	thoughtPrompt  string
	defaultPrompt  string
	functionPrompt string
	Thought        string
	Context        []services.Message
	ThoughtContext []services.Message
	service        services.ModelService
}

func NewAgent(service services.ModelService, agent *model.Agent) *Agent {

	result := &Agent{
		description:    agent,
		prompt:         agent.Prompt,
		service:        service,
		Context:        make([]services.Message, 0),
		ReplyChannel:   make(chan *messages.WebSocketMessage),
		CommandChannel: make(chan<- *messages.CommandType),
	}

	thoughtPrompt, err := embeddings.ReadFile("embedded/prompts/thought-prompt")
	defaultPrompt, err := embeddings.ReadFile("embedded/prompts/default")

	if err != nil {
		panic(err)
	}
	log.Println(fmt.Sprintf("Agent Prompt: %s", agent.Prompt))
	result.Context = append(result.Context, services.Message{Role: "system", Content: agent.Prompt})
	result.thoughtPrompt = string(thoughtPrompt)
	result.defaultPrompt = string(defaultPrompt)
	return result
}

func (agent *Agent) Chat(context string, text string) (err error) {

	agent.ReplyChannel <- messages.Message("thinking", messages.ReplyMessage{})
	agent.processThoughts(context, text)
	templateBuffer := bytes.NewBufferString("")
	defaultTemplate, err := template.New("default-prompt").Parse(agent.defaultPrompt)
	err = defaultTemplate.Execute(templateBuffer, map[string]string{"prompt": agent.prompt})
	rendered := templateBuffer.String()
	log.Println(fmt.Sprintf("Default Prompt: %s", rendered))
	agent.Context = append(agent.Context, services.Message{Role: "system", Content: rendered})

	summarizePrompt, err := embeddings.ReadFile("embedded/prompts/summarize-prompt")
	agent.Context = append(agent.Context, services.Message{Role: "system", Content: string(summarizePrompt)})
	agent.Context = append(agent.Context, agent.ThoughtContext[len(agent.ThoughtContext)-1])
	agent.Context = append(agent.Context, services.Message{Role: "user", Content: text})

	request := services.Request{Messages: agent.Context, Stream: false}
	err = agent.service.SendRequest(request, agent.responseHandler)
	if err != nil {
		return err
	}
	log.Println(agent.description.Name)
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

func (agent *Agent) responseHandler(message services.Message) {
	agent.Context = append(agent.Context, message)
	agent.ReplyChannel <- messages.Message("chat-response", messages.ReplyMessage{
		Agent:   agent.description.Name,
		Content: message.Content,
	})
}

func (agent *Agent) thoughtHandler(message services.Message) {

	agent.ThoughtContext = append(agent.ThoughtContext,
		services.Message{
			Role:    "assistant",
			Content: "<THOUGHT>" + message.Content + "</THOUGHT>",
		})
	agent.ReplyChannel <- messages.Message("thought-reply", messages.ReplyMessage{
		Agent:   agent.description.Name,
		Content: message.Content,
	})

}

func (agent *Agent) processThoughts(context string, text string) {
	reference := fetchInformation(context, text)
	thoughtValues := map[string]string{"memory": reference, "prompt": text}

	thoughtBuffer := bytes.NewBufferString("")
	thoughtTemplate, err := template.New("thought-prompt").Parse(agent.thoughtPrompt)
	err = thoughtTemplate.Execute(thoughtBuffer, thoughtValues)
	renderedThoughtPrompt := thoughtBuffer.String()
	log.Println(renderedThoughtPrompt)
	thoughtMessage := services.Message{Role: "system", Content: renderedThoughtPrompt}
	reprompt := services.Message{Role: "user", Content: text}

	thoughtRequest := services.Request{Messages: []services.Message{thoughtMessage, reprompt}, Stream: false}
	err = agent.service.SendRequest(thoughtRequest, agent.thoughtHandler)

	if err != nil {
		panic(err)
	}
}
