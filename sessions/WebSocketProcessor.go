package sessions

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"multivac.network/services/agents/providers"
)

type OutputProcessor[OUT any] interface {
	Output(...chan OUT)
}

type WebSocketProcessor struct {
	socket *websocket.Conn
}

func NewWebSocketProcessor(socket *websocket.Conn) *WebSocketProcessor {
	processor := &WebSocketProcessor{
		socket: socket,
	}
	return processor
}

func (processor *WebSocketProcessor) Output(channels ...chan *providers.Message) {
	for _, output := range channels {
		go processor.initialize(output)
	}
}

func (processor *WebSocketProcessor) initialize(output chan *providers.Message) {
	for {
		_, p, err := processor.socket.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		message := providers.Message{}
		err = json.Unmarshal(p, &message)
		if err != nil {
			log.Println(err)
			continue
		}
		output <- &message
	}
}
