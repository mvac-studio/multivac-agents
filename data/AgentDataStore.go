package data

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"multivac.network/services/agents/graph/model"
	"strings"
)

type AgentDataStore struct {
	collection *mongo.Collection
}

func NewAgentDataStore() *AgentDataStore {
	db := GetDatabase()
	return &AgentDataStore{
		collection: db.Collection("agents"),
	}

}

// CreateAgent Creates a new agent
func (store *AgentDataStore) CreateAgent(agent *model.Agent) (*model.Agent, error) {

	agent.Name = strings.ToLower(agent.Name)
	result, err := store.collection.InsertOne(context.Background(), agent)
	agent.ID = result.InsertedID.(primitive.ObjectID).Hex()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return agent, nil
}

// RetrieveAgents Retrieves all agents
func (store *AgentDataStore) RetrieveAgents() ([]*model.Agent, error) {
	cursor, err := store.collection.Find(context.Background(), bson.M{})
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

func (store *AgentDataStore) GetAgentsByIds(ids []string) ([]*AgentModel, error) {
	var results []*AgentModel
	oids := make([]primitive.ObjectID, 0)
	for _, id := range ids {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		oids = append(oids, oid)
	}
	filter := bson.M{"_id": bson.M{"$in": oids}}
	cursor, err := store.collection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	err = cursor.All(context.Background(), &results)
	return results, err
}

func (store *AgentDataStore) FindAgentById(id string) *model.Agent {

	var agent AgentModel
	oid, err := primitive.ObjectIDFromHex(id)
	err = store.collection.FindOne(context.Background(), bson.M{"_id": oid}).Decode(&agent)
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

func (store *AgentDataStore) FindAgent(name string) *model.Agent {

	var agent AgentModel
	err := store.collection.FindOne(context.Background(), bson.M{"key": strings.ToLower(name)}).Decode(&agent)
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
