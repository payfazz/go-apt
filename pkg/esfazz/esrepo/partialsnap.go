package esrepo

import (
	"github.com/payfazz/go-apt/pkg/esfazz"
	"github.com/payfazz/go-apt/pkg/esfazz/eventstore"
	"github.com/payfazz/go-apt/pkg/esfazz/snapstore"
)

// PartialSnapshotEventSourceRepository is event source repository which save snapshot every some number of event
func PartialSnapshotEventSourceRepository(
	eStore eventstore.EventStore,
	sStore snapstore.SnapshotStore,
	newAgg esfazz.AggregateFactory,
	eventCount int,
) EventSourceRepository {
	return &snapshotESRepository{
		eStore:     eStore,
		sStore:     sStore,
		newAgg:     newAgg,
		eventCount: eventCount,
	}
}
