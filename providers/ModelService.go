package providers

type ModelProvider interface {
	SendRequest(request Request, output chan Message) error
}
