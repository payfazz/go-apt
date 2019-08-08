package snapmongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateIdUniqueIndex create aggregate unique index
func CreateIdUniqueIndex(collection *mongo.Collection) error {
	index := mongo.IndexModel{
		Keys: bson.D{
			{"id", 1},
		},
		Options: options.Index().SetUnique(true),
	}
	_, err := collection.Indexes().CreateOne(context.Background(), index)
	return err
}
