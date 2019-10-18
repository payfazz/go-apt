package snapmongo

import (
	"context"
	"encoding/json"

	"github.com/payfazz/go-apt/pkg/esfazz/snapstore"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoSnapshotStore struct {
	collection *mongo.Collection
}

// Save is a function to save aggregate to database
func (s *mongoSnapshotStore) Save(ctx context.Context, id string, data json.RawMessage) error {
	dataMap := make(map[string]interface{})
	err := json.Unmarshal(data, &dataMap)
	if err != nil {
		return err
	}

	dataMap["id"] = id
	opt := options.Replace().SetUpsert(true)
	_, err = s.collection.ReplaceOne(ctx, bson.D{{Key: "id", Value: id}}, dataMap, opt)
	return err
}

// Find find aggregate in database based on id
func (s *mongoSnapshotStore) Find(ctx context.Context, id string) (json.RawMessage, error) {
	dataMap := make(map[string]interface{})
	err := s.collection.FindOne(ctx, bson.D{{Key: "id", Value: id}}).Decode(&dataMap)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	data, err := json.Marshal(dataMap)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// SnapshotStore is a constructor for MongoDB based snapshot store
func SnapshotStore(collection *mongo.Collection) snapstore.SnapshotStore {
	return &mongoSnapshotStore{
		collection: collection,
	}
}
