package eventmongo

import (
	"context"
	"encoding/json"
	"github.com/payfazz/go-apt/pkg/esfazz"
	"github.com/payfazz/go-apt/pkg/esfazz/eventstore"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

type mongoEventStore struct {
	collection *mongo.Collection
	indexOnce  sync.Once
}

// Save is function to save event to collection
func (m *mongoEventStore) Save(ctx context.Context, events ...*esfazz.EventPayload) ([]*esfazz.Event, error) {
	els := make([]interface{}, len(events))
	results := make([]*esfazz.Event, len(events))
	for i, event := range events {

		jsonRaw, err := json.Marshal(event.Data)
		if err != nil {
			return nil, err
		}

		data := make(map[string]interface{})
		err = json.Unmarshal(jsonRaw, &data)
		if err != nil {
			return nil, err
		}
		agg := esfazz.BaseAggregate{
			Id:      event.Aggregate.GetId(),
			Version: event.Aggregate.GetVersion(),
		}

		els[i] = eventLog{
			Type:      event.Type,
			Aggregate: agg,
			Data:      data,
		}
		results[i] = &esfazz.Event{
			Type:      event.Type,
			Aggregate: &agg,
			Data:      jsonRaw,
		}
	}

	_, err := m.collection.InsertMany(ctx, els)
	if err != nil {
		return nil, err
	}
	return results, nil
}

// FindNotApplied return function not applied to the aggregate
func (m *mongoEventStore) FindNotApplied(ctx context.Context, agg esfazz.Aggregate) ([]*esfazz.Event, error) {
	var results []*esfazz.Event
	filter := bson.D{
		{"aggregate.id", agg.GetId()},
		{"aggregate.version", bson.D{{"$gte", agg.GetVersion()}}},
	}
	opt := options.Find().SetSort(bson.D{{"aggregate.version", 1}})
	cur, err := m.collection.Find(ctx, filter, opt)

	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		elem := &eventLog{}
		err := cur.Decode(elem)
		if err != nil {
			return nil, err
		}

		data, err := json.Marshal(elem.Data)
		if err != nil {
			return nil, err
		}

		ev := &esfazz.Event{
			Type:      elem.Type,
			Aggregate: &elem.Aggregate,
			Data:      data,
		}
		results = append(results, ev)
	}
	return results, nil
}

// EventStore is constructor for event store using mongodb
func EventStore(collection *mongo.Collection) eventstore.EventStore {
	return &mongoEventStore{collection: collection}
}
