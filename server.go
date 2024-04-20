package main

import (
	"context"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"log"
	"multivac.network/services/agents/data"
	"multivac.network/services/agents/graph"
	"multivac.network/services/agents/processors"
	"multivac.network/services/agents/providers/groq"
	"net/http"
	"os"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	initializeData()

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	router := mux.NewRouter()
	router.HandleFunc("/chat/{group}/{jwt}", agentChat)
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func initializeData() {
	clientOptions := options.Client()

	clientOptions.ApplyURI("mongodb+srv://db-ngent-io.rcarmov.mongodb.net")
	clientOptions.SetRetryWrites(true)
	clientOptions.SetAppName("db-ngent-io")
	clientOptions.SetWriteConcern(writeconcern.Majority())
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "G6VuD^us",
	})

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}
	data.SetDatabase(client.Database("ngent"))
}

var contexts = make([]*processors.GroupProcessor, 0)

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

		groupStore := data.NewGroupDataStore()
		groupModel, err := groupStore.GetGroup(vars["group"])
		agentStore := data.NewAgentDataStore()
		agents, err := agentStore.GetAgentsByIds(groupModel.Agents)

		apikey := os.Getenv("GROQ_API_KEY")
		provider := groq.NewService("mixtral-8x7b-32768", apikey)

		socketInput := processors.NewSocketInputProcessor(ws)
		socketOutput := processors.NewSocketOutputProcessor(ws)
		group := processors.NewGroupProcessor(groupModel, provider)
		socketInput.ConversationOutput.To(group.Input)

		for _, agentModel := range agents {
			agent := processors.NewAgentProcessor(agentModel, provider)
			agent.To(socketOutput.AgentInput)
			err := group.AddAgent(agent)
			if err != nil {
				log.Println(err)
				continue
			}
		}

		contexts = append(contexts, group)
	}
}
