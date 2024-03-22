package agents

type Reply struct {
	Agent   string `json:"agent"`
	Content string `json:"content"`
	Thought string `json:"thought-prompt"`
}
