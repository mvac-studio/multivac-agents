package fireworks

import (
	"multivac.network/services/agents/providers"
)

// Service that implements the ModelProvider interface
type FireworksService struct {
	model  string
	apikey string
}

// NewService creates a new Ollama Service
func NewService(model string, apikey string) providers.ModelProvider {
	return &FireworksService{
		model:  model,
		apikey: apikey,
	}
}

// SendRequest implement the ModelProvider.SendRequest method for the Service type
func (s *FireworksService) SendRequest(request providers.Request) (*providers.Message, error) {

	return nil, nil
}
