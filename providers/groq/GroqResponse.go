package groq

type GroqFunctionCallDescriptor struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type GroqToolCall struct {
	Id       string                     `json:"id"`
	Type     string                     `json:"type"`
	Function GroqFunctionCallDescriptor `json:"function"`
}

type GroqResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role      string         `json:"role"`
			Content   string         `json:"content"`
			ToolCalls []GroqToolCall `json:"tool_calls"`
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
