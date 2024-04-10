package data

type AgentModel struct {
	ID          string `bson:"_id,omitempty"`
	Name        string `bson:"name"`
	Key         string `bson:"key"`
	Description string `bson:"description"`
	Engine      string `bson:"engine"`
	Prompt      string `bson:"prompt"`
}
