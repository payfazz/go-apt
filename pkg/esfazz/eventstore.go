package esfazz

import (
	"context"
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

// EventStore is an interface used for event store
type EventStore interface {
	Save(ctx context.Context, ev *EventPayload) (*EventLog, error)
	FindAfterAggregate(ctx context.Context, agg Aggregate) ([]*EventLog, error)
}

type postgresEventStore struct {
	tableName string
	model     *EventLog
}

// Save is a function to save event to event store
func (e *postgresEventStore) Save(ctx context.Context, ev *EventPayload) (*EventLog, error) {
	query, err := fazzdb.GetTransactionOrQueryContext(ctx)
	if err != nil {
		return nil, err
	}

	// if no aggregate, event will be related to new aggregate object
	if ev.Aggregate == nil {
		uuidV4, _ := uuid.NewV4()
		ev.Aggregate = &BaseAggregate{
			Id:      uuidV4.String(),
			Version: 0,
		}
	}

	dataJsonByte, err := json.Marshal(ev.Data)
	if err != nil {
		return nil, err
	}
	data := json.RawMessage(dataJsonByte)

	el := EventLogModel(e.tableName)
	el.EventType = ev.Type
	el.AggregateId = ev.Aggregate.GetId()
	el.AggregateVersion = ev.Aggregate.GetVersion()
	el.Data = data

	id, err := query.Use(el).InsertCtx(ctx, false)
	if err != nil {
		return nil, err
	}
	el.EventId = id.(int64)

	return el, nil
}

// FindAfterAggregate return all event that is not applied to the aggregate object
func (e *postgresEventStore) FindAfterAggregate(ctx context.Context, agg Aggregate) ([]*EventLog, error) {
	query, err := fazzdb.GetTransactionOrQueryContext(ctx)
	if err != nil {
		return nil, err
	}

	conditions := []fazzdb.SliceCondition{
		{Connector: fazzdb.CO_NONE, Field: "aggregate_id", Operator: fazzdb.OP_EQUALS, Value: agg.GetId()},
		{Connector: fazzdb.CO_AND, Field: "aggregate_version", Operator: fazzdb.OP_MORE_THAN_EQUALS, Value: agg.GetVersion()},
	}
	results, err := query.Use(e.model).
		WhereMany(conditions...).
		AllCtx(ctx)

	if err != nil {
		return nil, err
	}
	return results.([]*EventLog), err
}

// PostgresEventStore is a function to create new EventStore
func PostgresEventStore(tableName string) EventStore {
	return &postgresEventStore{
		tableName: tableName,
		model:     EventLogModel(tableName),
	}
}
