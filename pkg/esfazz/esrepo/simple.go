package esrepo

import (
	"github.com/payfazz/go-apt/pkg/esfazz"
	"github.com/payfazz/go-apt/pkg/esfazz/eventstore"
	"github.com/payfazz/go-apt/pkg/esfazz/snapstore"
)

// SimpleEventSourceRepository is simple event source repository without snapshot
func SimpleEventSourceRepository(store eventstore.EventStore, newAgg esfazz.AggregateFactory) EventSourceRepository {
	return &snapshotESRepository{
		eStore: store,
		sStore: snapstore.EmptySnapshotStore{},
		newAgg: newAgg,
	}
}
