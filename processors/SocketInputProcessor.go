package processors

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"log"
	"multivac.network/services/agents/messages"
)

type SocketInputProcessor struct {
	ConversationOutput *Output[*messages.ConversationMessage]
	socket             *websocket.Conn
}

// NewSocketInputProcessor creates a new socket processor
func NewSocketInputProcessor(socket *websocket.Conn) *SocketInputProcessor {

	processor := &SocketInputProcessor{
		ConversationOutput: NewOutputProcessor[*messages.ConversationMessage](),
		socket:             socket,
	}
	processor.initialize()
	return processor
}

func (sp *SocketInputProcessor) initialize() {
	go func() {
		for {
			_, p, err := sp.socket.ReadMessage()
			if err != nil {
				log.Println(err)
				break
			}
			message := messages.SocketMessage{}
			err = json.Unmarshal(p, &message)
			if err != nil {
				log.Println(err)
				continue
			}
			switch message.Type {
			case "chat-request":
				output := messages.ConversationMessage{}
				err := mapstructure.Decode(message.Content, &output)
				if err != nil {
					log.Println(err)
					continue
				}
				sp.ConversationOutput.output <- &output
			}
		}
	}()
}
