package eventmongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateAggregateUniqueIndex create aggregate unique index
func CreateAggregateUniqueIndex(collection *mongo.Collection) error {
	index := mongo.IndexModel{
		Keys: bson.D{
			{"aggregate.id", 1},
			{"aggregate.version", 1},
		},
		Options: options.Index().SetUnique(true),
	}
	_, err := collection.Indexes().CreateOne(context.Background(), index)
	return err
}
