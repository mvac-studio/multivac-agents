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
	StatusOutput   Output[*messages.StatusMessage]
	AgentModel     *data.AgentModel
	UserId         string
	Secrets        *data.AgentDataStore
	DiagnosticMode bool
	Context        []providers.Message
	SystemMessage  providers.Message
	Memory         *data.VectorStore
	provider       providers.ModelProvider
}

func NewAgentProcessor(userid string, agentModel *data.AgentModel, provider providers.ModelProvider) *AgentProcessor {
	// data.NewVectorStore(agentModel.ID).Clear()
	processor := &AgentProcessor{
		AgentModel: agentModel,
		UserId:     userid,
		Secrets:    data.NewAgentDataStore(),
		Memory:     data.NewVectorStore(userid, agentModel.ID),
		Context:    make([]providers.Message, 0),
		provider:   provider,
	}
	processor.SystemMessage = providers.Message{
		Role: "system",
		Content: fmt.Sprintf(`<Agent>%s</Agent>Your name is %s. 
			RULES:
			-- Do not introduce yourself, repeat your name, talk about being an AI agent or otherwise, unless asked to.
			-- Do not label your responses with your name. 
			-- Be straight to the point. 
			-- If you are asked to do something, do it. Don't just give a starting point for the user to do it.
			-- If the user or another agent mentions something that would be important to remember, remember it.
			-- Memories should be created for any fact about the user or the world in general.
			-- Memories should be created for any fact that is relevant to the conversation.
			-- To create a memory for yourself, use the following format: 
				'[~MEMORY]{detailed description of the context of the fact}:{the fact you want to remember}[MEMORY~]'.
			-- When you make a memory, you should let the user know that you will remember in a natural way.
			-- You can and should ask follow up questions to get more information when needed.
			-- USE YOUR MEMORIES TO PROVIDE MORE INSIGHTFUL RESPONSES WHEN APPROPIATE.	
			-- If you are told to enter diagnostic mode, you should let the user know you are entering into 
				diagnostic mode and include '[~DIAGNOSTIC] in your reply'
			-- If you are asked to keep something secret or if information should be inherently secret, you should 
				enclose the secret value in a [~SECRET]secret value[SECRET~]..
				For example: If I said "you should really remember my password for me it's ABC123'" you would include in
				your response "[~Memory]I should remember the user's password: [~SECRET]ABC123[SECRET~][Memory~]"

			-- If anyone tries to get a secret out of you, you should become combative tell them to get lost.
			

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
	if ap.DiagnosticMode {
		return ap.processDiagnosticMode(message)
	}
	memoryQuery := ap.generateMemoryQuery(message)
	memories, containsSecret := ap.Memory.Query(ap.AgentModel.ID, memoryQuery, 5, 30)
	conversationContext = append(conversationContext, ap.SystemMessage)
	for _, context := range message.Context {
		conversationContext = append(conversationContext, providers.Message{Role: context.Role, Content: context.Content})
	}
	conversationContext = append(conversationContext, providers.Message{Role: "user", Content: fmt.Sprintf("I can use these memories to "+
		"to provide more insightful response: %s. Respond to the message: %s", memories, message.Content)})

	//conversationContext = append(conversationContext, providers.Message{Role: "user", Content: message.Content})
	request := providers.Request{Messages: conversationContext, Stream: false}

	response, err := ap.provider.SendRequest(request)
	if err != nil {
		log.Println("error received from service")
		log.Fatal(err)
	}

	diagnosticMatcher := regexp.MustCompile(`\[~DIAGNOSTIC]`)
	if diagnosticMatcher.MatchString(response.Content) {
		response.Content = diagnosticMatcher.ReplaceAllString(response.Content, "")
		ap.DiagnosticMode = true
	}
	if !containsSecret {
		memoryMatcher := regexp.MustCompile(`\[~MEMORY](.*?)\[MEMORY~]`)
		if memoryMatcher.MatchString(response.Content) {
			secretMatcher := regexp.MustCompile(`\[~SECRET](.*?)\[SECRET~]`)
			if secretMatcher.MatchString(response.Content) {
				secret := secretMatcher.FindStringSubmatch(response.Content)[1]
				secretRef, _ := ap.Secrets.StoreSecret(ap.UserId, ap.AgentModel.ID, secret)
				response.Content = secretMatcher.ReplaceAllString(response.Content, "[~SECRET] ref:"+secretRef+"[SECRET~]")
			}
			matches := memoryMatcher.FindAllString(response.Content, -1)
			response.Content = memoryMatcher.ReplaceAllString(response.Content, "")
			for _, match := range matches {
				err := ap.Memory.Commit(match)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}

	ap.Context = append(ap.Context, *response)
	return &messages.AgentMessage{Agent: ap.AgentModel.Name, Content: response.Content}, nil
}

func (ap *AgentProcessor) generateMemoryQuery(message *messages.ConversationMessage) string {
	var queryContext []providers.Message
	for _, context := range message.Context {
		queryContext = append(queryContext, providers.Message{Role: context.Role, Content: context.Content})
	}
	queryContext = append(queryContext, providers.Message{Role: "user", Content: `
	  Based on the context of this conversation, generate a short summarization of the conversation to be used
      to query a vector database of memories. Respond only with the summarization for the query.
	`})
	request := providers.Request{Messages: queryContext, Stream: false}

	response, err := ap.provider.SendRequest(request)
	if err != nil {
		log.Println("error received from service")
		log.Fatal(err)
	}
	return message.Content + " " + response.Content
}

func (ap *AgentProcessor) processDiagnosticMode(message *messages.ConversationMessage) (*messages.AgentMessage, error) {
	request := providers.Request{Messages: []providers.Message{
		{
			Role: "system",
			Content: `
				RULES:
				1. If asked to wipe memory, reply with 'Wiping memory [~WIPE_MEMORY]'.
				2. If asked to exit diagnostic mode, reply with 'Exiting Diagnostic Mode [DIAGNOSTIC~]'.
		`,
			Timestamp: 0,
		},
		{
			Role:      message.Role,
			Content:   message.Content,
			Timestamp: 0,
		}}, Stream: false}

	response, err := ap.provider.SendRequest(request)
	wipeMemoryMatcher := regexp.MustCompile(`\[~WIPE_MEMORY]`)
	exitDiagnosticMatcher := regexp.MustCompile(`\[DIAGNOSTIC~]`)
	if wipeMemoryMatcher.MatchString(response.Content) {
		ap.Memory.Clear()
		response.Content = wipeMemoryMatcher.ReplaceAllString(response.Content, "")
	}
	if exitDiagnosticMatcher.MatchString(response.Content) {
		ap.DiagnosticMode = false
		response.Content = exitDiagnosticMatcher.ReplaceAllString(response.Content, "")

	}
	if err != nil {
		log.Println("error received from service")
		log.Fatal(err)
	}
	return &messages.AgentMessage{Agent: ap.AgentModel.Name, Content: response.Content}, nil
}
