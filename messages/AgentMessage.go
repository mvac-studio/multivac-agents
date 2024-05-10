package messages

type AgentMessage struct {
	StatusId string `json:"statusId"`
	Agent    string `json:"agent"`
	Content  string `json:"content"`
}

type StatusMessage struct {
	StatusId string `json:"statusId"`
	Agent    string `json:"agent"`
	Status   string `json:"status"`
}
