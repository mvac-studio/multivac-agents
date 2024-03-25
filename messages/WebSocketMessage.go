package messages

type WebSocketMessage struct {
	Type    string      `json:"type"`
	Created int64       `json:"created"`
	Payload interface{} `json:"payload"`
}
