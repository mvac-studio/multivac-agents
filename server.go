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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"multivac.network/services/agents/data"
	"multivac.network/services/agents/graph"
	"multivac.network/services/agents/processors"
	"multivac.network/services/agents/providers/groq"
	"multivac.network/services/agents/services/multivac-edges"
	"net/http"
	"os"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	client := initializeEdgeClient()
	initializeData(client)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	router := mux.NewRouter()
	router.HandleFunc("/chat/{group}/{jwt}", agentChat)
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func initializeEdgeClient() edges.EdgeServiceClient {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.Dial("multivac-edges-service.default.svc.cluster.local:50051", opts...)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	return edges.NewEdgeServiceClient(conn)
}

func initializeData(edgesService edges.EdgeServiceClient) {
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
	data.SetEdgesService(edgesService)
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
	valid, user := validUser(vars["jwt"])
	if valid {
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
		agents, err := agentStore.GetAgentsByGroup(context.Background(), groupModel.ID)

		apiKey := os.Getenv("GROQ_API_KEY")
		model := "llama3-70b-8192"

		provider := groq.NewService(model, apiKey)

		//provider := fireworks.NewService(model, apiKey, 3000)
		//apiKey := os.Getenv("FIREWORKS_API_KEY")
		//model := "llama-v3-70b-instruct-hf"

		socketInput := processors.NewSocketInputProcessor(ws)
		socketOutput := processors.NewSocketOutputProcessor(ws)
		group := processors.NewGroupProcessor(user.Name, user.UserID, groupModel, provider)
		socketInput.ConversationOutput.To(group.Input)
		group.FinalOutput.To(socketOutput.AgentInput)

		for _, agentModel := range agents {
			agent := processors.NewAgentProcessor(user.UserID, agentModel, provider)
			err := group.AddAgent(agent)
			if err != nil {
				log.Println(err)
				continue
			}
		}

		contexts = append(contexts, group)
	}
}
