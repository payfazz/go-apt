package esrepo

import (
	"context"
	"encoding/json"
	"github.com/payfazz/go-apt/pkg/esfazz"
	"github.com/payfazz/go-apt/pkg/esfazz/eventstore"
	"github.com/payfazz/go-apt/pkg/esfazz/snapstore"
)

type snapshotESRepository struct {
	eStore     eventstore.EventStore
	sStore     snapstore.SnapshotStore
	newAgg     esfazz.AggregateFactory
	eventCount int
}

// Save save event to event and snapshot store
func (s *snapshotESRepository) Save(ctx context.Context, event *esfazz.Event) error {
	err := s.eStore.Save(ctx, event)
	if err != nil {
		return err
	}

	err = s.saveSnapshot(ctx, event.Aggregate.GetId())
	return err
}

// Find find aggregate from snapshot and apply new event to this aggregate
func (s *snapshotESRepository) Find(ctx context.Context, id string) (esfazz.Aggregate, error) {
	agg, err := s.findSnapshot(ctx, id)
	if err != nil {
		return nil, err
	}

	evs, err := s.eStore.FindNotApplied(ctx, agg)
	if err != nil {
		return nil, err
	}

	err = agg.Apply(evs...)
	if err != nil {
		return nil, err
	}

	return agg, nil
}

func (s *snapshotESRepository) findSnapshot(ctx context.Context, id string) (esfazz.Aggregate, error) {
	agg := s.newAgg(id)
	rawData, err := s.sStore.Find(ctx, id)
	if err != nil {
		return nil, err
	}

	if rawData != nil {
		err = json.Unmarshal(rawData, agg)
		if err != nil {
			return nil, err
		}
	}
	return agg, nil
}

func (s *snapshotESRepository) saveSnapshot(ctx context.Context, id string) error {
	agg, err := s.findSnapshot(ctx, id)
	if err != nil {
		return err
	}

	evs, err := s.eStore.FindNotApplied(ctx, agg)
	if err != nil {
		return err
	}

	// don't save snapshot if below maximum not applied event count
	if len(evs) <= s.eventCount {
		return nil
	}

	err = agg.Apply(evs...)
	if err != nil {
		return err
	}

	data, err := json.Marshal(agg)
	if err != nil {
		return err
	}
	err = s.sStore.Save(ctx, id, data)
	return err
}

// SnapshotEventSourceRepository is event source repository which save snapshot every event saved
func SnapshotEventSourceRepository(
	eStore eventstore.EventStore,
	sStore snapstore.SnapshotStore,
	newAgg esfazz.AggregateFactory,
) EventSourceRepository {
	return &snapshotESRepository{
		eStore:     eStore,
		sStore:     sStore,
		newAgg:     newAgg,
		eventCount: 0,
	}
}
