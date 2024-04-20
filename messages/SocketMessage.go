package messages

// SocketMessage is a struct that represents a request message
type SocketMessage struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}
