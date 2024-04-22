package messages

type ConversationMessage struct {
	Context   []*ConversationMessage `json:"context"`
	Role      string                 `json:"role"`
	Content   string                 `json:"content"`
	Timestamp int64                  `json:"timestamp"`
}
