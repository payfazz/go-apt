package esrepo

import (
	"context"
	"encoding/json"
	"github.com/payfazz/go-apt/pkg/esfazz"
	"github.com/payfazz/go-apt/pkg/esfazz/eventstore"
	"github.com/payfazz/go-apt/pkg/esfazz/snapstore"
)

type eventSourceRepository struct {
	eventStore    eventstore.EventStore
	snapshotStore snapstore.SnapshotStore
	newAggregate  esfazz.AggregateFactory
	listeners     []EventListener
}

// Save save event to event and snapshot store
func (s *eventSourceRepository) Save(ctx context.Context, events ...*esfazz.EventPayload) error {
	eventLogs, err := s.eventStore.Save(ctx, events...)
	if err != nil {
		return err
	}

	for _, listener := range s.listeners {
		err = listener(ctx, eventLogs...)
		if err != nil {
			return err
		}
	}

	return nil
}

// Find find aggregate from snapshot and apply new event to this aggregate
func (s *eventSourceRepository) Find(ctx context.Context, id string) (esfazz.Aggregate, error) {
	agg := s.newAggregate(id)
	agg, err := s.findSnapshot(ctx, agg)
	if err != nil {
		return nil, err
	}

	evs, err := s.eventStore.FindNotApplied(ctx, agg)
	if err != nil {
		return nil, err
	}

	err = agg.Apply(evs...)
	if err != nil {
		return nil, err
	}

	return agg, nil
}

func (s *eventSourceRepository) findSnapshot(ctx context.Context, agg esfazz.Aggregate) (esfazz.Aggregate, error) {
	if s.snapshotStore == nil {
		return agg, nil
	}

	rawData, err := s.snapshotStore.Find(ctx, agg.GetId())
	if err != nil {
		return nil, err
	}

	if rawData != nil {
		agg = s.newAggregate(agg.GetId())
		err = json.Unmarshal(rawData, agg)
		if err != nil {
			return nil, err
		}
	}
	return agg, nil
}

func (s *eventSourceRepository) saveSnapshot(ctx context.Context, id string) error {
	if s.snapshotStore == nil {
		return nil
	}

	agg := s.newAggregate(id)
	agg, err := s.findSnapshot(ctx, agg)
	if err != nil {
		return err
	}

	evs, err := s.eventStore.FindNotApplied(ctx, agg)
	if err != nil {
		return err
	}

	err = agg.Apply(evs...)
	if err != nil {
		return err
	}

	data, err := json.Marshal(agg)
	if err != nil {
		return err
	}
	err = s.snapshotStore.Save(ctx, id, data)
	return err
}

// Build is constructor for event source repository
func Build(config Config) Repository {
	return &eventSourceRepository{
		eventStore:    config.eventStore,
		snapshotStore: config.snapshotStore,
		newAggregate:  config.aggregateFactory,
		listeners:     config.listeners,
	}
}
