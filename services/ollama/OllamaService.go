package ollama

import (
	"belli/services"
	"context"
	"time"
)
import "github.com/jmorganca/ollama/api"

// Service that implements the ModelService interface
type Service struct {
	client *api.Client
	model  string
}

// NewService creates a new Ollama Service
func NewService(model string) services.ModelService {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		panic(err)
	}
	return &Service{
		model:  model,
		client: client,
	}
}

// SendRequest implement the ModelService.SendRequest method for the Service type
func (s *Service) SendRequest(request services.Request, handler func(response services.Message)) error {
	chatRequest := &api.ChatRequest{Model: s.model, Messages: make([]api.Message, 0), Stream: &request.Stream}
	for _, message := range request.Messages {
		chatRequest.Messages = append(chatRequest.Messages, api.Message{Role: message.Role, Content: message.Content})
	}
	err := s.client.Chat(context.Background(), chatRequest, s.responseHandler(handler))
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) responseHandler(handler func(response services.Message)) func(response api.ChatResponse) error {
	return func(response api.ChatResponse) error {
		handler(services.Message{Role: response.Message.Role, Content: response.Message.Content,
			Timestamp: time.Now().UnixMilli(),
		})
		return nil
	}
}
