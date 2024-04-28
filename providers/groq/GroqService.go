package groq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"multivac.network/services/agents/providers"
	"multivac.network/services/agents/tools"
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
	groqRequest.ToolChoice = "none"
	if !request.DisableToolUse {
		groqRequest.AddTool(Tool{
			Type: "function",
			Function: FunctionDescriptor{
				Name:        "browse_web",
				Description: "Used to browse the webpage of an address and extract content and links.",
				Parameters: ParameterDescriptor{
					Type: "object",
					Properties: map[string]PropertyDescriptor{
						"address": {
							Type:        "string",
							Description: "fully qualified web address to visit use https://{address}"},
					},
					Required: []string{"address"}},
			}})
		groqRequest.AddTool(Tool{
			Type: "function",
			Function: FunctionDescriptor{
				Name:        "get_current_date",
				Description: "Should always be used to get the current date and or time.",
				Parameters: ParameterDescriptor{
					Type: "object",
					Properties: map[string]PropertyDescriptor{
						"format": {
							Type:        "string",
							Description: "the format of the date-time to use in golang time package format."},
					},
					Required: []string{"format"}},
			}})
	}
	for _, message := range request.Messages {
		groqRequest.Messages = append(groqRequest.Messages, GroqMessage{Role: message.Role, Content: message.Content})
	}
	log.Println(groqRequest)
	groqResponse, err := s.Send(groqRequest, client)

	if err != nil {
		return nil, err
	}
	if groqResponse.Choices[0].FinishReason == "tool_calls" {
		// execute tool call, add result to messages and send again
		for _, toolCall := range groqResponse.Choices[0].Message.ToolCalls {
			// execute tool call
			if toolCall.Function.Name == "browse_web" {
				var arguments map[string]string
				err := json.Unmarshal([]byte(toolCall.Function.Arguments), &arguments)
				if err != nil {
					return nil, err
				}
				result := tools.OpenWebAddress(arguments["address"])
				groqRequest.Messages = append(groqRequest.Messages, GroqMessage{Role: "assistant", Content: result})
				groqRequest.ToolChoice = "none"
				groqResponse, err = s.Send(groqRequest, client)
			}
			if toolCall.Function.Name == "get_current_date" {
				var arguments map[string]string
				err := json.Unmarshal([]byte(toolCall.Function.Arguments), &arguments)
				if err != nil {
					return nil, err
				}
				result := tools.GetCurrentDate(arguments["format"])
				groqRequest.Messages = append(groqRequest.Messages, GroqMessage{Role: "user", Content: result})
				groqRequest.ToolChoice = "none"
				groqResponse, err = s.Send(groqRequest, client)
			}
		}
	}
	return &providers.Message{
		Role:      groqResponse.Choices[0].Message.Role,
		Content:   groqResponse.Choices[0].Message.Content,
		Timestamp: time.Now().UnixMilli(),
	}, nil

}

func (s *Service) Send(groqRequest GroqRequest, client *http.Client) (*GroqResponse, error) {
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
	if resp.StatusCode != 200 {
		log.Println(fmt.Sprintf("Error: %s --V\n%s", resp.Status, body))
		return nil, err
	}
	log.Println(resp.Body)
	err = json.Unmarshal(body, &groqResponse)
	return &groqResponse, err
}
