package esfazz

import (
	"context"
	"encoding/json"
	"reflect"
)

// EventRepository is interface for event aggregate repository
type EventRepository interface {
	Save(ctx context.Context, payload *EventPayload) (Aggregate, error)
	Find(ctx context.Context, id string) (Aggregate, error)
}

type postgresEventRepository struct {
	aggregate  Aggregate
	eventStore EventStore
	aggStore   AggregateStore
}

// Save save account event and aggregate snapshot to storage
func (r *postgresEventRepository) Save(ctx context.Context, payload *EventPayload) (Aggregate, error) {
	// save to event store
	savedEvent, err := r.eventStore.Save(ctx, payload)
	if err != nil {
		return nil, err
	}

	agg, err := r.saveSnapshot(ctx, savedEvent.AggregateId)
	if err != nil {
		return nil, err
	}

	return agg, nil
}

// Find return account aggregate by id
func (r *postgresEventRepository) Find(ctx context.Context, id string) (Aggregate, error) {
	agg := reflect.New(reflect.ValueOf(r.aggregate).Elem().Type()).Interface().(Aggregate)
	agg.SetId(id)

	// load data from saved aggregate
	aggRow, err := r.aggStore.Find(ctx, id)
	if err != nil {
		return nil, err
	}

	if aggRow != nil {
		err = json.Unmarshal(aggRow.Data, agg)
		if err != nil {
			return nil, err
		}
	}

	// load new event and apply
	evs, err := r.eventStore.FindAfterAggregate(ctx, agg)
	if err != nil {
		return nil, err
	}

	for _, ev := range evs {
		err := agg.Apply(ev)
		if err != nil {
			return nil, err
		}
	}

	if agg.GetVersion() == 0 {
		return nil, nil
	}

	return agg, nil
}

func (r *postgresEventRepository) saveSnapshot(ctx context.Context, id string) (Aggregate, error) {
	agg, err := r.Find(ctx, id)
	if err != nil {
		return nil, err
	}

	_, err = r.aggStore.Save(ctx, agg)
	if err != nil {
		return nil, err
	}
	return agg, nil
}

// NewEventRepository create new event repository
func NewEventRepository(aggregate Aggregate, eventTable string, aggregateTable string) EventRepository {
	return &postgresEventRepository{
		aggregate:  aggregate,
		eventStore: PostgresEventStore(eventTable),
		aggStore:   PostgresAggregateStore(aggregateTable),
	}
}
