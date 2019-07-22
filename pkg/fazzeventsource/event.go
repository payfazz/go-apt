package fazzeventsource

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx/types"
	"time"
)

// EventLog is a struct for event
type EventLog struct {
	EventId          int64           `json:"event_id" db:"event_id"`
	EventType        string          `json:"event_type" db:"event_type"`
	AggregateId      string          `json:"aggregate_id" db:"aggregate_id"`
	AggregateVersion int             `json:"aggregate_version" db:"aggregate_version"`
	Data             json.RawMessage `json:"data" db:"data"`
	CreatedAt        *time.Time      `json:"created_at" db:"created_at"`
}

type EventPayload struct {
	Type             string
	AggregateId      string
	AggregateVersion int
	Data             interface{}
}

// EventStore is an interface used for event store
type EventStore interface {
	Save(ctx context.Context, ev EventPayload) (*EventLog, error)
	FindAllBy(ctx context.Context, aggregateId string, firstVersion int) ([]*EventLog, error)
}

type postgresEventStore struct {
	tableName string
}

// Save is a function to save event to event store
func (e *postgresEventStore) Save(ctx context.Context, ev EventPayload) (*EventLog, error) {

	dataJsonByte, err := json.Marshal(ev.Data)
	if err != nil {
		return nil, err
	}
	dataJsonText := types.JSONText(dataJsonByte)

	query, err := getContext(ctx)
	if err != nil {
		return nil, err
	}

	el := &EventLog{}
	queryGet := fmt.Sprintf(`INSERT INTO %s (event_type, aggregate_id, aggregate_version, data, created_at) 
									VALUES ($1,$2,$3,$4,$5) RETURNING *`, e.tableName)
	result, err := query.RawFirstCtx(
		ctx, el, queryGet, ev.Type, ev.AggregateId, ev.AggregateVersion, dataJsonText, time.Now())
	if err != nil {
		return nil, err
	}
	return result.(*EventLog), err
}

func (e *postgresEventStore) FindAllBy(ctx context.Context, aggregateId string, firstVersion int) ([]*EventLog, error) {
	query, err := getContext(ctx)
	if err != nil {
		return nil, err
	}
	el := &EventLog{}
	querySelect := fmt.Sprintf(`SELECT * FROM %s WHERE aggregate_id = $1 AND aggregate_version >= $2 ORDER BY event_id ASC`, e.tableName)
	results, err := query.RawAllCtx(ctx, el, querySelect, aggregateId, firstVersion)
	if err != nil {
		return nil, err
	}
	return results.([]*EventLog), err
}

// PostgresEventStore is a function to create new EventStore
func PostgresEventStore(tableName string) EventStore {
	return &postgresEventStore{
		tableName: tableName,
	}
}
