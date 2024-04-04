package groq

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"multivac.network/services/agents/services"
	"net/http"
	"time"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GroqRequest struct {
	Messages []Message `json:"messages"`
	Model    string    `json:"model"`
}

type GroqResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		Logprobs     any    `json:"logprobs"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int     `json:"prompt_tokens"`
		PromptTime       float64 `json:"prompt_time"`
		CompletionTokens int     `json:"completion_tokens"`
		CompletionTime   float64 `json:"completion_time"`
		TotalTokens      int     `json:"total_tokens"`
		TotalTime        float64 `json:"total_time"`
	} `json:"usage"`
	SystemFingerprint any `json:"system_fingerprint"`
}

// Service that implements the ModelService interface
type Service struct {
	model  string
	apikey string
}

// NewService creates a new Ollama Service
func NewService(model string, apikey string) services.ModelService {
	return &Service{
		model:  model,
		apikey: apikey,
	}
}

// SendRequest implement the ModelService.SendRequest method for the Service type
func (s *Service) SendRequest(request services.Request, handler func(response services.Message)) error {
	client := &http.Client{}
	groqRequest := GroqRequest{Messages: make([]Message, 0), Model: s.model}
	for _, message := range request.Messages {
		groqRequest.Messages = append(groqRequest.Messages, Message{Role: message.Role, Content: message.Content})
	}
	log.Println(groqRequest)
	data, err := json.Marshal(groqRequest)
	if err != nil {
		return err
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
		return err
	}
	handler(services.Message{
		Role:      groqResponse.Choices[0].Message.Role,
		Content:   groqResponse.Choices[0].Message.Content,
		Timestamp: time.Now().UnixMilli(),
	})
	return nil
}
