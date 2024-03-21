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
	prompt         string
	thoughtPrompt  string
	functionPrompt string
	Thought        string
	context        []services.Message
	service        services.ModelService
}

func NewAgent(service services.ModelService, agent *model.Agent) *Agent {

	result := &Agent{prompt: agent.Prompt, service: service, context: make([]services.Message, 0)}

	thoughtPrompt, err := embeddings.ReadFile("embedded/prompts/thought")
	functionPrompt, err := embeddings.ReadFile("embedded/prompts/function")

	if err != nil {
		panic(err)
	}
	println(agent.Prompt)
	result.context = append(result.context, services.Message{Role: "system", Content: agent.Prompt})
	result.thoughtPrompt = string(thoughtPrompt)
	result.functionPrompt = string(functionPrompt)
	return result
}

type Reply struct {
	Content string `json:"content"`
	Thought string `json:"thought"`
}

func (agent *Agent) Chat(text string) (reply Reply, err error) {

	reference := fetchInformation(text)
	fmt.Println(reference)
	referenceMessage := services.Message{Role: "assistant", Content: "USE THIS CONTENT TO ANSWER QUESTIONS <REF>" + reference + "</REF>"}
	message := services.Message{Role: "user", Content: text}

	thoughtBuffer := bytes.NewBufferString("")
	thoughtTemplate, err := template.New("thought").Parse(agent.thoughtPrompt)
	err = thoughtTemplate.Execute(thoughtBuffer, message)
	thoughtMessage := services.Message{Role: "assistant", Content: thoughtBuffer.String()}

	agent.context = append(agent.context, thoughtMessage)

	thoughtRequest := services.Request{Messages: agent.context, Stream: false}
	err = agent.service.SendRequest(thoughtRequest, agent.thoughtHandler)
	agent.context = append(agent.context, referenceMessage)
	if err != nil {
		panic(err)
	}
	agent.context = append(agent.context, message)
	request := services.Request{Messages: agent.context, Stream: false}
	err = agent.service.SendRequest(request, agent.responseHandler)
	if err != nil {
		panic(err)
	}

	return Reply{Content: agent.context[len(agent.context)-1].Content, Thought: agent.context[len(agent.context)-4].Content}, nil
}

func fetchInformation(text string) string {
	client := &http.Client{}
	query := url.QueryEscape(text)
	response, err := client.Get("http://multivac-embeddings-service.default.svc.cluster.local/search?q=" + query)
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
	agent.context = append(agent.context,
		services.Message{
			Role:    "assistant",
			Content: "<THOUGHT>" + message.Content + "</THOUGHT>",
		})
}
