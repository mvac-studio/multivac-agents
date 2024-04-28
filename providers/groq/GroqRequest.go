package groq

type PropertyDescriptor struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type ParameterDescriptor struct {
	Type       string                        `json:"type"`
	Properties map[string]PropertyDescriptor `json:"properties"`
	Required   []string                      `json:"required"`
}

type FunctionDescriptor struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Parameters  ParameterDescriptor `json:"parameters"`
}

type Tool struct {
	Type     string             `json:"type"`
	Function FunctionDescriptor `json:"function"`
}

type GroqRequest struct {
	Messages   []GroqMessage `json:"messages"`
	Model      string        `json:"model"`
	Tools      []Tool        `json:"tools"`
	ToolChoice string        `json:"tool_choice,default=auto"`
}

func (r *GroqRequest) AddTool(tool Tool) {
	if r.Tools == nil {
		r.Tools = make([]Tool, 0)
		r.ToolChoice = "auto"
	}
	r.Tools = append(r.Tools, tool)
}
