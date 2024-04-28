package providers

type Function struct {
	Description string            `json:"description"`
	Name        string            `json:"name"`
	Parameters  map[string]string `json:"parameters"`
}

type Tool struct {
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

type Request struct {
	Messages       []Message `json:"messages"`
	Stream         bool      `json:"stream"`
	Tools          []Tool    `json:"tools"`
	DisableToolUse bool      `json:"disable_tool_use"`
}
