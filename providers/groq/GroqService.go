package groq

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"multivac.network/services/agents/providers"
	"net/http"
	"time"
)

// Service that implements the ModelProvider interface
type Service struct {
	model  string
	apikey string
}

// NewService creates a new Ollama Service
func NewService(model string, apikey string) providers.ModelProvider {
	return &Service{
		model:  model,
		apikey: apikey,
	}
}

// SendRequest implement the ModelProvider.SendRequest method for the Service type
func (s *Service) SendRequest(request providers.Request) (*providers.Message, error) {
	client := &http.Client{}
	groqRequest := GroqRequest{Messages: make([]GroqMessage, 0), Model: s.model}
	for _, message := range request.Messages {
		groqRequest.Messages = append(groqRequest.Messages, GroqMessage{Role: message.Role, Content: message.Content})
	}
	log.Println(groqRequest)
	data, err := json.Marshal(groqRequest)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewReader(data))
	r.Header.Add("Authorization", "Bearer "+s.apikey)
	r.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(r)
	groqResponse := GroqResponse{}
	body, err := io.ReadAll(resp.Body)
	log.Println(resp.Body)
	err = json.Unmarshal(body, &groqResponse)
	if err != nil {
		return nil, err
	}
	return &providers.Message{
		Role:      groqResponse.Choices[0].Message.Role,
		Content:   groqResponse.Choices[0].Message.Content,
		Timestamp: time.Now().UnixMilli(),
	}, nil

}
