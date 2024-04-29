package providers

type Message struct {
	Role         string `json:"role"`
	Content      string `json:"content"`
	ImageContent string `json:"imageContent"`
	Timestamp    int64  `json:"timestamp"`
}
