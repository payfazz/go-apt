package fazzeventsource

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx/types"
	"github.com/payfazz/go-apt/pkg/fazzdb"
	"time"
)

// Event is a struct for event
type Event struct {
	Id        int64           `json:"id" db:"id"`
	Type      string          `json:"type" db:"type"`
	Data      json.RawMessage `json:"data" db:"data"`
	CreatedAt *time.Time      `json:"created_at" db:"created_at"`
}

// EventStore is an interface used for event store
type EventStore interface {
	Save(ctx context.Context, eventType string, eventData interface{}) (*Event, error)
	FindAllByKey(ctx context.Context, name string, value string, after int64) ([]*Event, error)
}

type postgresEventStore struct {
	tableName string
}

// Save is a function to save event to event store
func (e *postgresEventStore) Save(ctx context.Context, evType string, evData interface{}) (*Event, error) {

	dataJsonByte, err := json.Marshal(evData)
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
	result, err := query.RawFirstCtx(ctx, ev, queryGet, evType, dataJsonText, time.Now())
	if err != nil {
		return nil, err
	}
	return result.(*Event), err
}

func (e *postgresEventStore) FindAllByKey(ctx context.Context, name string, value string, after int64) ([]*Event, error) {
	query, err := fazzdb.GetQueryContext(ctx)
	if err != nil {
		return nil, err
	}
	ev := &Event{}
	querySelect := fmt.Sprintf(`SELECT * FROM %s WHERE data ->> '%s' = $1  AND id > $2 ORDER BY id ASC`, e.tableName, name)
	results, err := query.RawAllCtx(ctx, ev, querySelect, value, after)
	if err != nil {
		return nil, err
	}
	return results.([]*Event), err
}

// PostgresEventStore is a function to create new EventStore
func PostgresEventStore(tableName string) EventStore {
	return &postgresEventStore{
		tableName: tableName,
	}
}
