package providers

type Request struct {
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type ModelProvider interface {
	SendRequest(request Request, output chan Message) error
}
