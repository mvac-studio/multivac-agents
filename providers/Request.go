package providers

type Request struct {
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}
