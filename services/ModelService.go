package services

type Request struct {
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type Message struct {
	Role      string `json:"role"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

type ModelService interface {
	SendRequest(request Request, responseHandler func(response Message)) error
}
