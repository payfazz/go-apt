package eventpostgres

import (
	"context"
	"encoding/json"
	"github.com/payfazz/go-apt/pkg/esfazz"
	"github.com/payfazz/go-apt/pkg/esfazz/eventstore"
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

type postgresEventStore struct {
	tableName string
}

// Save is function to save event to database
func (p *postgresEventStore) Save(ctx context.Context, agg esfazz.Aggregate, events ...*esfazz.EventPayload) ([]*esfazz.Event, error) {
	query, err := fazzdb.GetTransactionOrQueryContext(ctx)
	if err != nil {
		return nil, err
	}

	models := make([]*eventLog, len(events))
	for i, ev := range events {
		dataRaw, err := json.Marshal(ev.Data)
		if err != nil {
			return nil, err
		}

		el := EventLogModel(p.tableName)
		el.EventType = ev.Type
		el.AggregateId = agg.GetId()
		el.AggregateVersion = agg.GetVersion() + int64(i)
		el.Data = dataRaw

		models[i] = el
	}
	_, err = query.Use(EventLogModel(p.tableName)).BulkInsertCtx(ctx, models)
	if err != nil {
		return nil, err
	}

	results := make([]*esfazz.Event, len(models))
	for i, ev := range models {
		results[i] = ev.ToEvent()
	}

	return results, nil
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
		OrderBy("aggregate_version", fazzdb.DIR_ASC).
		AllCtx(ctx)
	if err != nil {
		return nil, err
	}

	logs := queryRes.([]*eventLog)
	results := make([]*esfazz.Event, len(logs))
	for i, ev := range logs {
		results[i] = ev.ToEvent()
	}

	return results, nil
}

// EventStore is constructor for event store using postgres
func EventStore(tableName string) eventstore.EventStore {
	return &postgresEventStore{
		tableName: tableName,
	}
}
