package repository

import (
	"context"
	"github.com/payfazz/go-apt/example/esfazz/domain/account/command/aggregate"
	"github.com/payfazz/go-apt/pkg/esfazz"
	"github.com/payfazz/go-apt/pkg/esfazz/esrepo"
	"github.com/payfazz/go-apt/pkg/esfazz/eventstore"
	"github.com/payfazz/go-apt/pkg/esfazz/snapstore"
)

// AccountEventRepository is repository for account event
type AccountEventRepository interface {
	Save(ctx context.Context, payload *esfazz.Event) (*aggregate.Account, error)
	Find(ctx context.Context, id string) (*aggregate.Account, error)
}

type accountEventRepository struct {
	repository esrepo.EventSourceRepository
}

// Save save account event and aggregate snapshot to storage
func (a *accountEventRepository) Save(ctx context.Context, payload *esfazz.Event) (*aggregate.Account, error) {
	// save to event store
	err := a.repository.Save(ctx, payload)
	if err != nil {
		return nil, err
	}
	return a.Find(ctx, payload.Aggregate.GetId())
}

// Find return account aggregate by id
func (a *accountEventRepository) Find(ctx context.Context, id string) (*aggregate.Account, error) {
	result, err := a.repository.Find(ctx, id)
	if err != nil {
		return nil, err
	}
	return result.(*aggregate.Account), nil
}

// NewAccountEventRepository create new account event repository
func NewAccountEventRepository(eventStore eventstore.EventStore, snapStore snapstore.SnapshotStore) AccountEventRepository {

	return &accountEventRepository{
		repository: esrepo.SnapshotEventSourceRepository(eventStore, snapStore, aggregate.AccountAggregate),
	}
}
