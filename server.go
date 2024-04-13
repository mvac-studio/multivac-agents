package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"multivac.network/services/agents/data"
	"multivac.network/services/agents/executors"
	"multivac.network/services/agents/graph"
	"multivac.network/services/agents/providers/groq"
	"multivac.network/services/agents/sessions"
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
	router.HandleFunc("/chat/{group}/{jwt}", agentChat)
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

var chatSessions = make([]*sessions.GroupContext, 0)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func agentChat(writer http.ResponseWriter, request *http.Request) {
	// TODO: check authorization

	vars := mux.Vars(request)
	log.Println(vars["jwt"])
	log.Println(vars["group"])
	if validUser(vars["jwt"]) {
		//TODO (jkelly): validate the actual authorization header
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		log.Println("Upgrading connection:", request.RemoteAddr)
		ws, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			log.Println(err)
		}
		s := data.NewAgentStore()
		agentModel := s.FindAgent(vars["agent"])
		apikey := os.Getenv("GROQ_API_KEY")
		var agent = executors.NewAgent(groq.NewService("mixtral-8x7b-32768", apikey), agentModel)
		chatSessions = append(chatSessions, sessions.NewSession(vars["jwt"], ws, agent))
	}
}
