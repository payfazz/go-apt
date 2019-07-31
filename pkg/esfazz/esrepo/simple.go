package esrepo

import (
	"context"
	"github.com/payfazz/go-apt/pkg/esfazz"
	"github.com/payfazz/go-apt/pkg/esfazz/eventstore"
)

type simpleESRepository struct {
	store  eventstore.EventStore
	newAgg esfazz.AggregateFactory
}

// Save save event to repository
func (s *simpleESRepository) Save(ctx context.Context, event *esfazz.Event) error {
	return s.store.Save(ctx, event)
}

// Find find aggregate from repository by creating from saved event
func (s *simpleESRepository) Find(ctx context.Context, id string) (interface{}, error) {
	agg := s.newAgg(id)
	evs, err := s.store.FindNotApplied(ctx, agg)

	if err != nil {
		return nil, err
	}

	for _, ev := range evs {
		err := agg.Apply(ev)
		if err != nil {
			return nil, err
		}
	}

	return agg, nil
}

// SimpleEventSourceRepository is simple event source repository without snapshot
func SimpleEventSourceRepository(store eventstore.EventStore, newAgg esfazz.AggregateFactory) EventSourceRepository {
	return &simpleESRepository{
		store:  store,
		newAgg: newAgg,
	}
}
