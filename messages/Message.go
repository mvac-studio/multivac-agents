package messages

import "time"

type Message struct {
	Type    string      `json:"type"`
	Created int64       `json:"created"`
	Payload interface{} `json:"payload"`
}

func CreateMessage(t string, p interface{}) *Message {
	return &Message{Type: t, Created: time.Now().UnixMilli(), Payload: p}
}
