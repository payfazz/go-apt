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

	agg := &esfazz.BaseAggregate{
		Id:      "01234567-89ab-cdef-0123-456789abcdef",
		Version: 0,
	}

	evs, err := store.Save(ctx, agg, &esfazz.EventPayload{
		Type: "test.event",
		Data: map[string]string{"test": "234"},
	})

	if err != nil {
		t.Errorf("error saving data: %s", err)
	}

	if len(evs) == 0 {
		t.Error("save result empty")
	}
}

func TestMongoEventStore_FindNotApplied(t *testing.T) {
	ctx, collection := prepareContextAndCollection()
	store := EventStore(collection)
	agg := &esfazz.BaseAggregate{
		Id:      "01234567-89ab-cdef-0123-456789abcdef",
		Version: 0,
	}

	events := []*esfazz.EventPayload{
		{
			Type: "test.event",
			Data: map[string]string{"test": "234"},
		},
		{
			Type: "test.event",
			Data: map[string]string{"test": "345"},
		},
		{
			Type: "test.event",
			Data: map[string]string{"test": "456"},
		},
	}

	evs, err := store.Save(ctx, agg, events...)
	if err != nil {
		t.Errorf("error saving data: %s", err)
	}
	if len(evs) != 3 {
		t.Errorf("saved event list not of same length, expected: 3, result %d", len(evs))
	}

	evResults, err := store.FindNotApplied(ctx, agg)

	if err != nil {
		t.Errorf("error saving data: %s", err)
	}
	if len(evResults) != 3 {
		t.Errorf("data in list not on the same length as existing data, expected: 3, result: %d", len(evResults))
	}

	for idx := range events {
		if evResults[idx].Aggregate.GetVersion() != int64(idx) {
			t.Errorf("event in index %d is not in expected order", idx)
		}
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
