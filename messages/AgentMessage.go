package messages

type AgentMessage struct {
	Agent   string `json:"agent"`
	Content string `json:"content"`
}

type StatusMessage struct {
	Agent  string `json:"agent"`
	Status string `json:"status"`
}
