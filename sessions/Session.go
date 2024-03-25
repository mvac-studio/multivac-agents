package sessions

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"multivac.network/services/agents/executors"
	"multivac.network/services/agents/services"
)

type Session struct {
	UserId string
	socket *websocket.Conn
	Agent  *executors.Agent
}

func NewSession(userId string, socket *websocket.Conn, agent *executors.Agent) *Session {

	result := &Session{UserId: userId, socket: socket, Agent: agent}

	go result.start()
	return result
}

func (session *Session) start() {
	go func() {
		for {
			select {
			case message := <-session.Agent.ReplyChannel:

				output, err := json.Marshal(message)
				err = session.SendMessage(string(output))
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()
	for {
		messageType, p, err := session.socket.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(messageType, string(p))
		message := services.Message{}
		err = json.Unmarshal(p, &message)
		if err != nil {
			log.Println(err)
			return
		}
		err = session.Agent.Chat(session.UserId, message.Content)
		if err != nil {
			return
		}
	}
}

func (chatSession *Session) SendMessage(message string) error {
	println(message)
	return chatSession.socket.WriteMessage(websocket.TextMessage, []byte(message))
}

func (chatSession *Session) Close() error {
	return chatSession.socket.Close()
}
