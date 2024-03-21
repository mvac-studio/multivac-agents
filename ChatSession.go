package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"multivac.network/services/agents/agents"
	"multivac.network/services/agents/services"
)

type ChatSession struct {
	UserId string
	socket *websocket.Conn
	Agent  *agents.Agent
}

func NewChatSession(userId string, socket *websocket.Conn, agent *agents.Agent) *ChatSession {

	result := &ChatSession{UserId: userId, socket: socket, Agent: agent}

	go result.start()
	return result
}

func (session *ChatSession) start() {
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
		response, _ := session.Agent.Chat(session.UserId, message.Content)
		output, err := json.Marshal(response)
		log.Println("response: ", string(output))
		err = session.SendMessage(string(output))
		if err != nil {
			return
		}
	}
}

func (chatSession *ChatSession) SendMessage(message string) error {
	println(message)
	return chatSession.socket.WriteMessage(websocket.TextMessage, []byte(message))
}

func (chatSession *ChatSession) Close() error {
	return chatSession.socket.Close()
}
