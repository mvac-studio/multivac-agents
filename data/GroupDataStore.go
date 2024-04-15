package data

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

var database *mongo.Database

func SetDatabase(db *mongo.Database) {
	database = db
}

func GetDatabase() (db *mongo.Database) {
	return database
}

// GroupDataStore struct that uses a MongoDB collection
type GroupDataStore struct {
	collection *mongo.Collection
}

// NewGroupDataStore function creates a new GroupDataStore using a mongodb database parameter
func NewGroupDataStore() *GroupDataStore {
	db := GetDatabase()
	return &GroupDataStore{
		collection: db.Collection("groups"),
	}
}

// CreateGroup function creates a new group in the database
func (g *GroupDataStore) CreateGroup(group *GroupModel) (*GroupModel, error) {
	group.Agents = make([]string, 0)
	result, err := g.collection.InsertOne(context.Background(), group)
	if err != nil {
		log.Fatalln(err)
	}
	group.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return group, err
}

// GetGroup function retrieves a group from the database by ID
func (g *GroupDataStore) GetGroup(id string) (*GroupModel, error) {
	group := &GroupModel{}
	objectId, err := primitive.ObjectIDFromHex(id)
	err = g.collection.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(group)
	return group, err
}

// GetAllGroups function retrieves all groups regardless of owner from the database
func (g *GroupDataStore) GetAllGroups() ([]*GroupModel, error) {
	cursor, err := g.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(nil)

	var groups = make([]*GroupModel, 0)
	for cursor.Next(nil) {
		group := &GroupModel{}
		err := cursor.Decode(group)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	return groups, nil
}

// GetGroups function retrieves all groups from the database
func (g *GroupDataStore) GetGroups(owner string) ([]*GroupModel, error) {
	cursor, err := g.collection.Find(context.Background(), bson.M{"owner": owner})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(nil)

	var groups = make([]*GroupModel, 0)
	for cursor.Next(nil) {
		group := &GroupModel{}
		err := cursor.Decode(group)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	return groups, nil
}

// UpdateGroup function updates a group in the database
func (g *GroupDataStore) UpdateGroup(group *GroupModel) error {
	_, err := g.collection.ReplaceOne(nil, group.ID, group)
	return err
}

// DeleteGroup function deletes a group from the database by ID
func (g *GroupDataStore) DeleteGroup(id string) error {
	_, err := g.collection.DeleteOne(nil, id)
	return err
}

func (g *GroupDataStore) AddAgentToGroup(id string, agentID string) (*GroupModel, error) {
	bsonId, err := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": bsonId}
	update := bson.M{
		"$push": bson.M{"agents": agentID},
		"$set":  bson.M{"updated": int(time.Now().Unix())},
	}

	result := g.collection.FindOneAndUpdate(context.Background(), filter, update)

	group := &GroupModel{}
	err = result.Decode(group)
	return group, err
}

func (g *GroupDataStore) RemoveAgentFromGroup(id string, agentID string) error {
	_, err := g.collection.UpdateOne(nil, id, map[string]string{"$pull": agentID})
	return err
}

func (g *GroupDataStore) FindGroupsByAgentId(agentID string) ([]*GroupModel, error) {
	cursor, err := g.collection.Find(context.Background(), bson.M{"agents": agentID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(nil)

	var groups = make([]*GroupModel, 0)
	for cursor.Next(nil) {
		group := &GroupModel{}
		err := cursor.Decode(group)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}

	return groups, nil
}

// ArchiveGroup function archives a group in the database by ID
func (g *GroupDataStore) ArchiveGroup(id string) error {
	_, err := g.collection.UpdateOne(nil, id, map[string]bool{"archived": true})
	return err
}
