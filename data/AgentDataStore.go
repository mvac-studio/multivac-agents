package data

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"multivac.network/services/agents/graph/model"
	edges "multivac.network/services/agents/services/multivac-edges"
	"net/url"
	"strings"
)

type AgentDataStore struct {
	collection *mongo.Collection
	edges      edges.EdgeServiceClient
}

type Vertex struct {
	Ref  string `bson:"ref"`
	Type string `bson:"type"`
}

type EdgeModel struct {
	ID      string `bson:"_id,omitempty"`
	Target  Vertex `bson:"target"`
	Source  Vertex `bson:"source"`
	Created int64  `bson:"created"`
	Updated int64  `bson:"updated"`
}

func NewAgentDataStore() *AgentDataStore {
	db := GetDatabase()
	return &AgentDataStore{
		collection: db.Collection("agents"),
		edges:      edgesService,
	}

}

// DeleteAgent Deletes an agent
func (store *AgentDataStore) DeleteAgent(ctx context.Context, id string) (*AgentModel, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	deletedAgent := &AgentModel{}
	result := store.collection.FindOneAndDelete(ctx, bson.M{"_id": oid})
	err = result.Decode(deletedAgent)
	return deletedAgent, err

}

// SaveAgent Creates a new agent
func (store *AgentDataStore) SaveAgent(ctx context.Context, agent *model.Agent) (*model.Agent, error) {

	id, _ := primitive.ObjectIDFromHex(agent.ID)
	agent.Key = url.PathEscape(strings.ToLower(agent.Name))

	dataModel := AgentModel{
		Name:        agent.Name,
		Key:         agent.Key,
		Description: agent.Description,
		Engine:      agent.Engine,
		Prompt:      agent.Prompt,
	}

	if id != primitive.NilObjectID {
		filter := bson.M{"_id": id}
		opts := options.Update().SetUpsert(true)
		result, err := store.collection.UpdateOne(ctx, filter, bson.M{"$set": dataModel}, opts)

		if result.UpsertedID != nil {
			agent.ID = result.UpsertedID.(primitive.ObjectID).Hex()
		}
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
	} else {
		result, err := store.collection.InsertOne(ctx, dataModel)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		agent.ID = result.InsertedID.(primitive.ObjectID).Hex()
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

func (store *AgentDataStore) GetAgentsByGroup(ctx context.Context, groupid string) ([]*AgentModel, error) {
	result, err := store.edges.GetForwardEdges(ctx, &edges.GetForwardEdgesRequest{
		Source:     &edges.Vertex{Ref: groupid, Type: "group"},
		TargetType: "agent",
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	ids := []string{}
	for _, edge := range result.Edges {
		ids = append(ids, edge.Target.Ref)
	}
	return store.GetAgentsByIds(ids)
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

func (store *AgentDataStore) StoreSecret(userid string, agentid string, value string) (string, error) {
	db := GetDatabase()
	secretCollection := db.Collection("secrets")
	secret := MemorySecret{
		AgentID: agentid,
		UserID:  userid,
		Secret:  value,
	}

	op, err := secretCollection.InsertOne(context.Background(), secret)
	return op.InsertedID.(primitive.ObjectID).Hex(), err
}

func (store *AgentDataStore) RetrieveSecret(id string, userid string, agentid string) (string, error) {
	db := GetDatabase()
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	secretCollection := db.Collection("secrets")
	var secret MemorySecret
	err = secretCollection.FindOne(context.Background(), bson.M{"_id": oid, "agent_id": agentid, "user_id": userid}).Decode(&secret)
	if secret.Secret == "" {
		secret.Secret = "[Redacted]"
	}
	return secret.Secret, err
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

type MemorySecret struct {
	Id      string `bson:"_id,omitempty"`
	AgentID string `bson:"agent_id"`
	UserID  string `bson:"user_id"`
	Secret  string `bson:"secret"`
}
