package es

import (
	"context"
	"encoding/json"
	"github.com/jmoiron/sqlx/types"
	"github.com/payfazz/go-apt/pkg/fazzdb"
	"time"
)

// EventRepository is an interface used for event store
type EventRepository interface {
	Save(ctx context.Context, eventType string, eventData interface{}) (*Event, error)
	AggregateById(ctx context.Context, id string) ([]Event, error)
}

type eventRepository struct{}

// Save is a function to save event to event store
func (e *eventRepository) Save(ctx context.Context, eventType string, eventData interface{}) (*Event, error) {

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

// AggregateById return all events related to id in data sorted from last event
func (e *eventRepository) AggregateById(ctx context.Context, id string) ([]Event, error) {
	query, err := fazzdb.GetQueryContext(ctx)
	if err != nil {
		return nil, err
	}

	var evs []Event
	querySelect := `SELECT * FROM events WHERE data ->> 'id' = $1 ORDER BY id DESC`
	err = query.Db.SelectContext(ctx, &evs, querySelect, id)
	return evs, err
}

// NewEventStore is a function to create new EventRepository
func NewEventStore() EventRepository {
	return &eventRepository{}
}
