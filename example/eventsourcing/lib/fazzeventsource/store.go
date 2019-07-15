package fazzeventsource

import (
	"context"
	"encoding/json"
	"github.com/jmoiron/sqlx/types"
	"github.com/payfazz/go-apt/example/eventsourcing/lib/fazzpubsub"
	"github.com/payfazz/go-apt/pkg/fazzdb"
	"time"
)

// EventStore is an interface used for event store
type EventStore interface {
	Save(ctx context.Context, eventType string, eventData interface{}) (*Event, error)
	Publish(ctx context.Context, topic string, event *Event) error
	FindByInstanceId(ctx context.Context, id string) ([]*Event, error)
}

type postgresEventRepository struct {
	pubsub fazzpubsub.PubSub
}

// Save is a function to save event to event store
func (e *postgresEventRepository) Save(ctx context.Context, eventType string, eventData interface{}) (*Event, error) {

	dataJsonByte, err := json.Marshal(eventData)
	if err != nil {
		return nil, err
	}
	dataJsonText := types.JSONText(dataJsonByte)

	query, err := fazzdb.GetQueryContext(ctx)
	if err != nil {
		return nil, err
	}

	var ev = &Event{}
	queryGet := `INSERT INTO events (type,data,created_at) VALUES ($1,$2,$3) RETURNING *`
	result, err := query.RawFirstCtx(ctx, ev, queryGet, eventType, dataJsonText, time.Now())
	return result.(*Event), err
}

// Publish do publishing to event pubsub
func (e *postgresEventRepository) Publish(ctx context.Context, topic string, event *Event) error {
	evJson, err := json.Marshal(event)
	if err != nil {
		return err
	}
	err = e.pubsub.Publish(ctx, topic, evJson)
	return err

}

// FindByInstanceId return all events related to id in data sorted from last event
func (e *postgresEventRepository) FindByInstanceId(ctx context.Context, id string) ([]*Event, error) {
	query, err := fazzdb.GetQueryContext(ctx)
	if err != nil {
		return nil, err
	}

	var ev = &Event{}
	querySelect := `SELECT * FROM events WHERE data ->> 'id' = $1 ORDER BY id DESC`
	results, err := query.RawAllCtx(ctx, ev, querySelect, id)
	evs := results.([]*Event)
	return evs, err
}

// NewEventStore is a function to create new EventStore
func NewEventStore(pubsub fazzpubsub.PubSub) EventStore {
	return &postgresEventRepository{pubsub: pubsub}
}
