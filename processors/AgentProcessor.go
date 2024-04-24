package processors

import (
	"fmt"
	"log"
	"multivac.network/services/agents/data"
	"multivac.network/services/agents/messages"
	"multivac.network/services/agents/providers"
	"regexp"
)

type AgentProcessor struct {
	Processor[*messages.ConversationMessage, *messages.AgentMessage]
	StatusOutput Output[*messages.StatusMessage]
	AgentModel   *data.AgentModel
	Context      []providers.Message
	Memory       *data.VectorStore
	provider     providers.ModelProvider
}

func NewAgentProcessor(agentModel *data.AgentModel, provider providers.ModelProvider) *AgentProcessor {
	// data.NewVectorStore(agentModel.ID).Clear()
	processor := &AgentProcessor{
		AgentModel: agentModel,
		Memory:     data.NewVectorStore(agentModel.ID),
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
			"marked with <Agent>agentname</Agent>. %s. IF EITHER THE REQUEST OR YOUR RESPONSE SHOULD BE REMEMBERED LONG "+
			"RESPOND ALSO WITH '<MEMORY>{detailed description of the context of the memory} {the memory you want to "+
			"remember}</MEMORY>'. You should only remember information that is important to not forget. Don't tell the "+
			"user you are updating your memory.", agentModel.Name, agentModel.Name, agentModel.Prompt),
	})
	processor.Processor = NewProcessor[*messages.ConversationMessage, *messages.AgentMessage](processor.Process)
	return processor
}

func (ap *AgentProcessor) Process(message *messages.ConversationMessage) (*messages.AgentMessage, error) {
	var conversationContext []providers.Message
	for _, context := range message.Context {
		conversationContext = append(conversationContext, providers.Message{Role: context.Role, Content: context.Content})
	}
	conversationContext = append(conversationContext, providers.Message{Role: "assistant", Content: fmt.Sprintf("This is what I remember, I can use this memory to "+
		"to provide more insightful response. <MEMORY>%s</MEMORY>", ap.Memory.Query(message.Content, 5, 30))})

	conversationContext = append(conversationContext, providers.Message{Role: "user", Content: message.Content})
	request := providers.Request{Messages: conversationContext, Stream: false}

	response, err := ap.provider.SendRequest(request)
	if err != nil {
		log.Println("error received from service")
		log.Fatal(err)
	}
	matcher := regexp.MustCompile(`<MEMORY>(.*?)</MEMORY>`)
	if matcher.MatchString(response.Content) {
		matches := matcher.FindAllString(response.Content, -1)
		response.Content = matcher.ReplaceAllString(response.Content, "")
		for _, match := range matches {
			err := ap.Memory.Commit(match)
			if err != nil {
				log.Println(err)
			}
		}
	}

	ap.Context = append(ap.Context, *response)
	return &messages.AgentMessage{Agent: ap.AgentModel.Name, Content: response.Content}, nil
}
