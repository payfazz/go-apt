package eventmongo

import (
	"context"
	_ "github.com/lib/pq"
	"github.com/payfazz/go-apt/pkg/esfazz"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"testing"
)

func TestMongoEventStore_Save(t *testing.T) {
	ctx, collection := prepareContextAndCollection()

	store := EventStore(collection)

	err := store.Save(ctx, &esfazz.Event{
		Type: "test.event",
		Aggregate: &esfazz.BaseAggregate{
			Id:      "01234567-89ab-cdef-0123-456789abcdef",
			Version: 0,
		},
		Data: []byte("{\"test\": \"234\"}"),
	})

	if err != nil {
		t.Errorf("error saving data: %s", err)
	}
}

func TestMongoEventStore_FindNotApplied(t *testing.T) {
	ctx, collection := prepareContextAndCollection()

	store := EventStore(collection)
	err := store.Save(ctx, &esfazz.Event{
		Type: "test.event",
		Aggregate: &esfazz.BaseAggregate{
			Id:      "01234567-89ab-cdef-0123-456789abcdef",
			Version: 0,
		},
		Data: []byte("{\"test\": \"234\"}"),
	})
	if err != nil {
		t.Errorf("error saving data: %s", err)
	}

	err = store.Save(ctx, &esfazz.Event{
		Type: "test.event",
		Aggregate: &esfazz.BaseAggregate{
			Id:      "01234567-89ab-cdef-0123-456789abcdef",
			Version: 1,
		},
		Data: []byte("{\"test\": \"345\"}"),
	})
	if err != nil {
		t.Errorf("error saving data: %s", err)
	}

	err = store.Save(ctx, &esfazz.Event{
		Type: "test.event",
		Aggregate: &esfazz.BaseAggregate{
			Id:      "01234567-89ab-cdef-0123-456789abcdef",
			Version: 2,
		},
		Data: []byte("{\"test\": \"456\"}"),
	})
	if err != nil {
		t.Errorf("error saving data: %s", err)
	}

	evs, err := store.FindNotApplied(ctx,
		&esfazz.BaseAggregate{
			Id:      "01234567-89ab-cdef-0123-456789abcdef",
			Version: 0,
		},
	)
	if err != nil {
		t.Errorf("error saving data: %s", err)
	}
	if len(evs) != 3 {
		t.Errorf("existing data not list not on the same length, expected: 3, result: %d", len(evs))
	}
}

func prepareContextAndCollection() (context.Context, *mongo.Collection) {
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

	collection := client.Database("test").Collection("events")
	_ = collection.Drop(ctx)

	return ctx, collection
}
