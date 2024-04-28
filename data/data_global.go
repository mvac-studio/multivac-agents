package data

import (
	"go.mongodb.org/mongo-driver/mongo"
	"multivac.network/services/agents/services/multivac-edges"
)

var database *mongo.Database
var edgesService edges.EdgeServiceClient

func SetDatabase(db *mongo.Database) {
	database = db
}
