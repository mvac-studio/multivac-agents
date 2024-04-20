package providers

type ModelProvider interface {
	SendRequest(request Request) (*Message, error)
}
