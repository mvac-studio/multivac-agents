package sessions

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"multivac.network/services/agents/providers"
)

type GroupContext struct {
	group    string
	socket   *websocket.Conn
	model    providers.ModelProvider
	Input    chan providers.Message
	internal chan providers.Message
	Output   chan providers.Message
}

func NewGroupContext(group string, socket *websocket.Conn, provider providers.ModelProvider) *GroupContext {
	result := &GroupContext{group: group, socket: socket, model: provider}
	go result.start()
	return result
}

func (group *GroupContext) start() {
	go group.initializeInternal()
	go group.initializeOutput()
	go group.initializeInput()
	go group.initializeSocket()
}

func (group *GroupContext) initializeInput() {
	for {
		select {
		case message := <-group.Input:
			request := providers.Request{Messages: make([]providers.Message, 0), Stream: false}

			// TODO: create a new message with the group context, agent list, and instructions to delegate
			// request.Messages = append(request.Messages, message)

			request.Messages = append(request.Messages, message)
			err := group.model.SendRequest(request, group.internal)
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}
}

func (group *GroupContext) initializeOutput() {
	for {
		select {
		case message := <-group.Output:
			output, err := json.Marshal(message)
			err = group.emitMessage(string(output))
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (group *GroupContext) emitMessage(message string) error {
	return group.socket.WriteMessage(websocket.TextMessage, []byte(message))
}

func (group *GroupContext) Close() error {
	return group.socket.Close()
}

func (group *GroupContext) initializeSocket() {
	for {
		_, p, err := group.socket.ReadMessage()
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
		group.Input <- message
	}
}

func (group *GroupContext) initializeInternal() {
	for {
		select {
		case message := <-group.internal:
			// TODO: Instead parse and dispatch to the agents in the internal response.
			group.Output <- message
		}
	}
}
