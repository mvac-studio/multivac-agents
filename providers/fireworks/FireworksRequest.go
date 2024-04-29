package fireworks

type FireworksResponseFormatContainer struct {
	Type FireworksResponseFormat `json:"type,omitempty"`
}

type FireworksRequest struct {
	Messages             []FireworksRequestMessage        `json:"messages,omitempty"`
	MaxTokens            int                              `json:"max_tokens,omitempty"`
	PromptTruncateLength int                              `json:"prompt_truncate_len,omitempty"`
	Temperature          float64                          `json:"temperature,omitempty"`
	TopP                 float64                          `json:"top_p,omitempty"`
	FrequencyPenalty     float64                          `json:"frequency_penalty,omitempty"`
	PresencePenalty      float64                          `json:"presence_penalty,omitempty"`
	N                    int                              `json:"n,omitempty"`
	Stop                 []string                         `json:"stop,omitempty"`
	ResponseFormat       FireworksResponseFormatContainer `json:"response_format,omitempty"`
	Stream               bool                             `json:"stream,omitempty"`
	Model                string                           `json:"model,omitempty"`
}
