package processors

import (
	"github.com/gorilla/websocket"
	"multivac.network/services/agents/messages"
)

type SocketOutputProcessor struct {
	AgentInput  *Input[*messages.AgentMessage]
	StatusInput *Input[*messages.StatusMessage]
	socket      *websocket.Conn
}

// NewSocketOutputProcessor creates a new socket processor
func NewSocketOutputProcessor(socket *websocket.Conn) *SocketOutputProcessor {

	processor := &SocketOutputProcessor{
		AgentInput:  NewInputProcessor[*messages.AgentMessage](),
		StatusInput: NewInputProcessor[*messages.StatusMessage](),
		socket:      socket,
	}
	processor.initialize()
	return processor
}

func (sp *SocketOutputProcessor) initialize() {
	go func() {
		for {
			response := messages.SocketMessage{Type: "chat-response"}
			response.Content = <-sp.AgentInput.input
			err := sp.socket.WriteJSON(response)
			if err != nil {
				break
			}
		}
	}()
	go func() {
		for {
			response := messages.SocketMessage{Type: "status-event"}
			response.Content = <-sp.StatusInput.input
			err := sp.socket.WriteJSON(response)
			if err != nil {
				break
			}
		}
	}()
}
