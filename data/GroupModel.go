package data

// GroupModel is a bson serializable data model.
type GroupModel struct {
	ID          string   `bson:"_id,omitempty"`
	Name        string   `bson:"name"`
	Owner       string   `bson:"owner"`
	Description string   `bson:"description"`
	Created     int      `bson:"created"`
	Updated     int      `bson:"updated"`
	Archived    bool     `bson:"archived"`
	Agents      []string `bson:"agents"`
}
