package fazzeventsource

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx/types"
	"github.com/payfazz/go-apt/pkg/fazzdb"
	"time"
)

// EventStore is an interface used for event store
type EventStore interface {
	Save(ctx context.Context, eventType string, eventData interface{}) (*Event, error)
	FindByInstanceId(ctx context.Context, id string) ([]*Event, error)
}

type postgresEventStore struct {
	tableName string
}

// Save is a function to save event to event store
func (e *postgresEventStore) Save(ctx context.Context, eventType string, eventData interface{}) (*Event, error) {

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
	queryGet := fmt.Sprintf(`INSERT INTO %s (type,data,created_at) VALUES ($1,$2,$3) RETURNING *`, e.tableName)
	result, err := query.RawFirstCtx(ctx, ev, queryGet, eventType, dataJsonText, time.Now())
	return result.(*Event), err
}

// FindByInstanceId return all events related to id in data sorted from last event
func (e *postgresEventStore) FindByInstanceId(ctx context.Context, id string) ([]*Event, error) {
	query, err := fazzdb.GetQueryContext(ctx)
	if err != nil {
		return nil, err
	}

	var ev = &Event{}
	querySelect := fmt.Sprintf(`SELECT * FROM %s WHERE data ->> 'id' = $1 ORDER BY id ASC`, e.tableName)
	results, err := query.RawAllCtx(ctx, ev, querySelect, id)
	evs := results.([]*Event)
	return evs, err
}

// NewPostgresEventStore is a function to create new EventStore
func NewPostgresEventStore(tableName string) EventStore {
	return &postgresEventStore{
		tableName: tableName,
	}
}
