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
		Role:    "system",
		Content: fmt.Sprintf("Your name is %s. %s", agentModel.Name, agentModel.Prompt),
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
