package agents

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"multivac.network/services/agents/graph/model"
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
	prompt         string
	thoughtPrompt  string
	defaultPrompt  string
	functionPrompt string
	Thought        string
	context        []services.Message
	service        services.ModelService
}

func NewAgent(service services.ModelService, agent *model.Agent) *Agent {

	result := &Agent{description: agent, prompt: agent.Prompt, service: service, context: make([]services.Message, 0)}

	thoughtPrompt, err := embeddings.ReadFile("embedded/prompts/thought-prompt")
	defaultPrompt, err := embeddings.ReadFile("embedded/prompts/default")

	if err != nil {
		panic(err)
	}
	log.Println(fmt.Sprintf("Agent Prompt: %s", agent.Prompt))
	result.context = append(result.context, services.Message{Role: "system", Content: agent.Prompt})
	result.thoughtPrompt = string(thoughtPrompt)
	result.defaultPrompt = string(defaultPrompt)
	return result
}

func (agent *Agent) Chat(context string, text string) (reply Reply, err error) {

	agent.processThoughts(context, text)

	templateBuffer := bytes.NewBufferString("")
	defaultTemplate, err := template.New("default-prompt").Parse(agent.defaultPrompt)
	err = defaultTemplate.Execute(templateBuffer, map[string]string{"prompt": agent.prompt})
	rendered := templateBuffer.String()
	log.Println(fmt.Sprintf("Default Prompt: %s", rendered))
	agent.context = append(agent.context, services.Message{Role: "system", Content: rendered})

	summarizePrompt, err := embeddings.ReadFile("embedded/prompts/summarize-prompt")
	agent.context = append(agent.context, services.Message{Role: "user", Content: string(summarizePrompt)})
	request := services.Request{Messages: agent.context, Stream: false}
	err = agent.service.SendRequest(request, agent.responseHandler)
	if err != nil {
		panic(err)
	}
	log.Println(agent.description.Name)
	return Reply{Agent: agent.description.Name, Content: agent.context[len(agent.context)-1].Content, Thought: agent.context[len(agent.context)-2].Content}, nil
}

func fetchInformation(index string, text string) string {
	client := &http.Client{}
	query := url.QueryEscape(text)
	response, err := client.Get("http://multivac-embeddings-service.default.svc.cluster.local/search?index=" + index + "&q=" + query)
	if err != nil {
		log.Println(err)
		return ""
	}
	content, err := io.ReadAll(response.Body)
	return string(content)
}

func (agent *Agent) responseHandler(message services.Message) {
	agent.context = append(agent.context, message)
}

func (agent *Agent) thoughtHandler(message services.Message) {
	log.Println(fmt.Sprintf("Thought: %s", message.Content))
	agent.context = append(agent.context,
		services.Message{
			Role:    "assistant",
			Content: "<THOUGHT>" + message.Content + "</THOUGHT>",
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
