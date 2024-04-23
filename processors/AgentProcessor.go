package processors

import (
	"fmt"
	"log"
	"multivac.network/services/agents/data"
	"multivac.network/services/agents/messages"
	"multivac.network/services/agents/providers"
)

type AgentProcessor struct {
	Processor[*messages.ConversationMessage, *messages.AgentMessage]
	StatusOutput Output[*messages.StatusMessage]
	AgentModel   *data.AgentModel
	Context      []providers.Message
	provider     providers.ModelProvider
}

func NewAgentProcessor(agentModel *data.AgentModel, provider providers.ModelProvider) *AgentProcessor {
	processor := &AgentProcessor{
		AgentModel: agentModel,
		Context:    make([]providers.Message, 0),
		provider:   provider,
	}
	processor.Context = append(processor.Context, providers.Message{
		Role: "system",
		Content: fmt.Sprintf("<Agent>%s</Agent>Your name is %s. Do not introduce yourself, repeat your name, "+
			"talk about being an AI agent or otherwise, unless asked to. "+
			"Do not label your responses with your name. Be straight to the point. If you are asked to do something. "+
			"Do it, don't give a starting point for doing it. Complete the task. Other agents and their prompts and "+
			"responses will be "+
			"marked with <Agent>agentname</Agent>. %s", agentModel.Name, agentModel.Name, agentModel.Prompt),
	})
	processor.Processor = NewProcessor[*messages.ConversationMessage, *messages.AgentMessage](processor.Process)
	return processor
}

func (ap *AgentProcessor) Process(message *messages.ConversationMessage) (*messages.AgentMessage, error) {

	for _, context := range message.Context {
		ap.Context = append(ap.Context, providers.Message{Role: context.Role, Content: context.Content})
	}

	request := providers.Request{Messages: ap.Context, Stream: false}

	response, err := ap.provider.SendRequest(request)
	if err != nil {
		log.Println("error received from service")
		log.Fatal(err)
	}
	ap.Context = append(ap.Context, *response)
	return &messages.AgentMessage{Agent: ap.AgentModel.Name, Content: response.Content}, nil
}
