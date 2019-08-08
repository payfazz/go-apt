package esrepo

import (
	"context"
	"encoding/json"
	"github.com/payfazz/go-apt/pkg/esfazz"
	"github.com/payfazz/go-apt/pkg/esfazz/eventstore"
	"github.com/payfazz/go-apt/pkg/esfazz/snapstore"
)

type snapshotESRepository struct {
	eStore       eventstore.EventStore
	sStore       snapstore.SnapshotStore
	newAgg       esfazz.AggregateFactory
	snapEventMax int
}

// Save save event to event and snapshot store
func (s *snapshotESRepository) Save(ctx context.Context, event *esfazz.EventPayload) error {
	err := s.eStore.Save(ctx, event)
	if err != nil {
		return err
	}

	err = s.saveSnapshot(ctx, event.Aggregate.GetId())
	return err
}

// Find find aggregate from snapshot and apply new event to this aggregate
func (s *snapshotESRepository) Find(ctx context.Context, id string) (esfazz.Aggregate, error) {
	agg := s.newAgg(id)
	agg, err := s.findSnapshot(ctx, agg)
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

func (s *snapshotESRepository) findSnapshot(ctx context.Context, agg esfazz.Aggregate) (esfazz.Aggregate, error) {
	if s.sStore == nil {
		return agg, nil
	}

	rawData, err := s.sStore.Find(ctx, agg.GetId())
	if err != nil {
		return nil, err
	}

	if rawData != nil {
		agg = s.newAgg(agg.GetId())
		err = json.Unmarshal(rawData, agg)
		if err != nil {
			return nil, err
		}
	}
	return agg, nil
}

func (s *snapshotESRepository) saveSnapshot(ctx context.Context, id string) error {
	if s.sStore == nil {
		return nil
	}

	agg := s.newAgg(id)
	agg, err := s.findSnapshot(ctx, agg)
	if err != nil {
		return err
	}

	evs, err := s.eStore.FindNotApplied(ctx, agg)
	if err != nil {
		return err
	}

	//  save snapshot only if above maximum buffered event
	if len(evs) <= s.snapEventMax {
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

// snapshotEventSourceRepository is event source repository which save snapshot every event saved
func snapshotEventSourceRepository(
	eStore eventstore.EventStore,
	sStore snapstore.SnapshotStore,
	newAgg esfazz.AggregateFactory,
	eventMax int,
) Repository {
	return &snapshotESRepository{
		eStore:       eStore,
		sStore:       sStore,
		newAgg:       newAgg,
		snapEventMax: eventMax,
	}
}
