package espostgres

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jmoiron/sqlx/types"
	"github.com/payfazz/go-apt/pkg/esfazz"
	"github.com/payfazz/go-apt/pkg/esfazz/eventstore"
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

type postgresEventStore struct {
	tableName string
}

// Save is function to save event to database
func (p *postgresEventStore) Save(ctx context.Context, event *esfazz.Event) error {
	query, err := fazzdb.GetTransactionOrQueryContext(ctx)
	if err != nil {
		return err
	}

	if event.Aggregate.GetId() == "" {
		return errors.New("aggregate id for event must not be empty")
	}

	el := EventLogModel(p.tableName)
	el.EventType = event.Type
	el.AggregateId = event.Aggregate.GetId()
	el.AggregateVersion = event.Aggregate.GetVersion()
	el.Data = types.JSONText(event.Data)

	_, err = query.Use(el).InsertCtx(ctx, false)
	return err
}

// FindNotApplied return function not applied to the aggregate
func (p *postgresEventStore) FindNotApplied(ctx context.Context, agg esfazz.Aggregate) ([]*esfazz.Event, error) {
	query, err := fazzdb.GetTransactionOrQueryContext(ctx)
	if err != nil {
		return nil, err
	}

	conditions := []fazzdb.SliceCondition{
		{Connector: fazzdb.CO_NONE, Field: "aggregate_id", Operator: fazzdb.OP_EQUALS, Value: agg.GetId()},
		{Connector: fazzdb.CO_AND, Field: "aggregate_version", Operator: fazzdb.OP_MORE_THAN_EQUALS, Value: agg.GetVersion()},
	}

	queryRes, err := query.Use(EventLogModel(p.tableName)).
		WhereMany(conditions...).
		AllCtx(ctx)
	if err != nil {
		return nil, err
	}

	logs := queryRes.([]*EventLog)
	results := make([]*esfazz.Event, len(logs))
	for i, v := range logs {
		results[i] = &esfazz.Event{
			Type: v.EventType,
			Aggregate: &esfazz.BaseAggregate{
				Id:      v.AggregateId,
				Version: v.AggregateVersion,
			},
			Data: json.RawMessage(v.Data),
		}
	}

	return results, nil
}

// EventStore is constructor for event store using postgres
func EventStore(tableName string) eventstore.EventStore {
	return &postgresEventStore{
		tableName: tableName,
	}
}
