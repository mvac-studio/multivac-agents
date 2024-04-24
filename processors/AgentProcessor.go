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
	StatusOutput  Output[*messages.StatusMessage]
	AgentModel    *data.AgentModel
	Context       []providers.Message
	SystemMessage providers.Message
	Memory        *data.VectorStore
	provider      providers.ModelProvider
}

func NewAgentProcessor(agentModel *data.AgentModel, provider providers.ModelProvider) *AgentProcessor {
	// data.NewVectorStore(agentModel.ID).Clear()
	processor := &AgentProcessor{
		AgentModel: agentModel,
		Memory:     data.NewVectorStore(agentModel.ID),
		Context:    make([]providers.Message, 0),
		provider:   provider,
	}
	processor.SystemMessage = providers.Message{
		Role: "system",
		Content: fmt.Sprintf(`<Agent>%s</Agent>Your name is %s. 
			RULES:
			1. Do not introduce yourself, repeat your name, talk about being an AI agent or otherwise, unless asked to.
			2. Do not label your responses with your name. 
			3. Be straight to the point. 
			4. If you are asked to do something, do it. Don't just give a starting point for the user to do it.
			5. If the user or another agent mentions something that would be important to remember, remember it.
			6. To create a memory for yourself to remember, use the following format: 
			'[~MEMORY]{detailed description of the context of the memory} {the memory you want to remember}[MEMORY~]'.
			7. You can and should ask follow up questions to get more information when needed.
			

			INFORMATION:
			1. Other agents and their messages and responses will be marked with <Agent>agentname</Agent>. 

			ABOUT YOU:
			%s

			`, agentModel.Name, agentModel.Name, agentModel.Prompt),
	}
	processor.Processor = NewProcessor[*messages.ConversationMessage, *messages.AgentMessage](processor.Process)
	return processor
}

func (ap *AgentProcessor) Process(message *messages.ConversationMessage) (*messages.AgentMessage, error) {
	var conversationContext []providers.Message
	conversationContext = append(conversationContext, ap.SystemMessage)
	for _, context := range message.Context {
		conversationContext = append(conversationContext, providers.Message{Role: context.Role, Content: context.Content})
	}
	conversationContext = append(conversationContext, providers.Message{Role: "assistant", Content: fmt.Sprintf("This is what I remember, I can use this memory to "+
		"to provide more insightful response. %s", ap.Memory.Query(message.Content, 5, 30))})

	conversationContext = append(conversationContext, providers.Message{Role: "user", Content: message.Content})
	request := providers.Request{Messages: conversationContext, Stream: false}

	response, err := ap.provider.SendRequest(request)
	if err != nil {
		log.Println("error received from service")
		log.Fatal(err)
	}
	matcher := regexp.MustCompile(`\[~MEMORY](.*?)\[MEMORY~]`)
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
