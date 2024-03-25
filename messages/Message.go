package messages

import "time"

func Message(t string, p interface{}) *WebSocketMessage {
	return &WebSocketMessage{Type: t, Created: time.Now().UnixMilli(), Payload: p}
}
