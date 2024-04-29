package fireworks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"multivac.network/services/agents/providers"
	"net/http"
	"strings"
)

// Service that implements the ModelProvider interface
type FireworksService struct {
	model     string
	apikey    string
	maxTokens int
}

// NewService creates a new Ollama Service
func NewService(model string, apikey string, maxTokens int) providers.ModelProvider {
	return &FireworksService{
		model:     model,
		apikey:    apikey,
		maxTokens: maxTokens,
	}
}

type FireworksChoice struct {
	Index        int                      `json:"index,omitempty"`
	Message      FireworksResponseMessage `json:"message,omitempty"`
	FinishReason string                   `json:"finish_reason,omitempty"`
}

type FireworksUsage struct {
	PromptTokens     int `json:"prompt_tokens,omitempty"`
	TotalTokens      int `json:"total_tokens,omitempty"`
	CompletionTokens int `json:"completion_tokens,omitempty"`
}

type FireworksResponse struct {
	Id      string            `json:"id"`
	Object  string            `json:"object"`
	Created int64             `json:"created"`
	Model   string            `json:"model"`
	Choices []FireworksChoice `json:"choices"`
	Usage   FireworksUsage    `json:"usage"`
}

// SendRequest implement the ModelProvider.SendRequest method for the Service type
func (s *FireworksService) SendRequest(request providers.Request) (*providers.Message, error) {

	fireworksRequest := s.requestToFireworksRequest(request)
	data, err := json.Marshal(fireworksRequest)
	r, err := http.NewRequest("POST", "https://api.fireworks.ai/inference/v1/chat/completions", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	r.Header.Add("accept", "application/json")
	r.Header.Add("content-type", "application/json")
	r.Header.Add("authorization", "Bearer akbRouOJOuzXc3oRQfsnp4tALvTkIVyezYzDgmUPvcA5WVGH")

	response, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}
	requestBody, err := io.ReadAll(r.Body)
	fmt.Println(string(requestBody))
	body, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}
	var fireworksResponse FireworksResponse
	err = json.Unmarshal(body, &fireworksResponse)

	if err != nil {
		return nil, err

	}
	if len(fireworksResponse.Choices) == 0 {
		return nil, nil
	}
	return &providers.Message{
		Role:      fireworksResponse.Choices[0].Message.Role,
		Content:   fireworksResponse.Choices[0].Message.Content.Text,
		Timestamp: fireworksResponse.Created,
	}, nil
	return nil, nil
}

func (s *FireworksService) requestToFireworksRequest(request providers.Request) *FireworksRequest {
	return &FireworksRequest{
		Messages:             messagesToFireworksMessages(request.Messages),
		MaxTokens:            s.maxTokens,
		ResponseFormat:       FireworksResponseFormatContainer{Type: ResponseFormatText},
		PromptTruncateLength: 7000,
		Temperature:          0.8,
		TopP:                 1,
		FrequencyPenalty:     0.0,
		PresencePenalty:      0.0,
		N:                    1,
		// Stop:                 []string{},
		Stream: false,
		Model:  fmt.Sprintf("accounts/fireworks/models/%s", s.model),
	}
}

func messagesToFireworksMessages(messages []providers.Message) []FireworksRequestMessage {
	var fireworksMessages []FireworksRequestMessage
	fireworksContents := []FireworksContent{}
	for _, message := range messages {

		if message.ImageContent == "" {
			fireworksContents = append(fireworksContents, FireworksContent{
				Text: strings.Replace(strings.Replace(message.Content, "\n", " ", -1), "\t", " ", -1),
				Type: "text",
			})
		} else {
			fireworksContents = append(fireworksContents, FireworksContent{
				ImageUrl: &FireworksImageUrl{Url: message.ImageContent},
				Type:     "image_url",
			})
		}
	}
	fireworksMessages = append(fireworksMessages, FireworksRequestMessage{
		Role:    messages[0].Role,
		Content: fireworksContents,
	})
	return fireworksMessages
}
