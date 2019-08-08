package config

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
)

var mc *mongo.Client
var once sync.Once

// GetMongoClient return singleton mongo client
func GetMongoClient() *mongo.Client {
	once.Do(func() {
		ctx := context.Background()
		clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Fatal(err)
		}

		err = client.Ping(ctx, nil)
		if err != nil {
			log.Fatal(err)
		}

		mc = client
	})
	return mc
}
