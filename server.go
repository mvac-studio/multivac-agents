package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"multivac.network/services/agents/agents"
	"multivac.network/services/agents/graph"
	"multivac.network/services/agents/services/groq"
	"multivac.network/services/agents/store"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	router := mux.NewRouter()
	router.Handle("/api", playground.Handler("GraphQL playground", "/api/query"))
	router.Handle("/api/query", srv)
	router.HandleFunc("/chat/{agent}/{jwt}", agentChat)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

var sessions = make([]*ChatSession, 0)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func agentChat(writer http.ResponseWriter, request *http.Request) {
	// TODO: check authorization

	vars := mux.Vars(request)
	log.Println(vars["jwt"])
	log.Println(vars["agent"])
	if validUser(vars["jwt"]) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		ws, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			log.Println(err)
		}
		s := store.NewAgentStore()
		agentModel := s.FindAgent(vars["agent"])
		var agent = agents.NewAgent(groq.NewService("mixtral-8x7b-32768", os.Getenv("GROQ_API_KEY")), agentModel)
		sessions = append(sessions, NewChatSession(vars["jwt"], ws, agent))
	}
}
