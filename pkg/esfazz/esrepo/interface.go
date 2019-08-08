package esrepo

import (
	"context"
	"github.com/payfazz/go-apt/pkg/esfazz"
)

// Repository is interface for event repository
type Repository interface {
	Save(ctx context.Context, event *esfazz.EventPayload) error
	Find(ctx context.Context, id string) (esfazz.Aggregate, error)
}

// NewRepository is constructor for repository
func NewRepository(config *RepositoryConfig) Repository {
	return snapshotEventSourceRepository(
		config.EventStore,
		config.SnapshotStore,
		config.AggregateFactory,
		config.SnapshotEventMax,
	)
}
