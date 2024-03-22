package store

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"multivac.network/services/agents/graph/model"
	"strings"
)

type AgentStore struct {
	client *mongo.Client
}

func NewAgentStore() *AgentStore {
	clientOptions := options.Client().ApplyURI("mongodb://192.168.88.209:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return &AgentStore{client: client}
}

// CreateAgent Creates a new agent
func (store *AgentStore) CreateAgent(agent *model.Agent) (*model.Agent, error) {
	collection := store.client.Database("vector").Collection("agents")
	agent.Name = strings.ToLower(agent.Name)
	_, err := collection.InsertOne(context.Background(), agent)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return agent, nil
}

// RetrieveAgents Retrieves all agents
func (store *AgentStore) RetrieveAgents() ([]*model.Agent, error) {
	collection := store.client.Database("vector").Collection("agents")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer cursor.Close(context.Background())

	var results []*model.Agent

	for cursor.Next(context.Background()) {
		var agent AgentModel
		err := cursor.Decode(&agent)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		results = append(results, &model.Agent{
			ID:          agent.ID,
			Name:        agent.Name,
			Key:         agent.Key,
			Description: agent.Description,
			Engine:      agent.Engine,
			Prompt:      agent.Prompt,
		})
	}
	return results, nil
}

func (store *AgentStore) FindAgent(name string) *model.Agent {
	collection := store.client.Database("vector").Collection("agents")
	var agent AgentModel
	err := collection.FindOne(context.Background(), bson.M{"key": strings.ToLower(name)}).Decode(&agent)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return &model.Agent{
		ID:          agent.ID,
		Name:        agent.Name,
		Key:         agent.Key,
		Description: agent.Description,
		Engine:      agent.Engine,
		Prompt:      agent.Prompt,
	}

}

type AgentModel struct {
	ID          string `bson:"_id"`
	Name        string `bson:"name"`
	Key         string `bson:"key"`
	Description string `bson:"description"`
	Engine      string `bson:"engine"`
	Prompt      string `bson:"prompt"`
}
