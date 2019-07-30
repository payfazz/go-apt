package repository

import (
	"context"
	"github.com/payfazz/go-apt/example/esfazz/domain/account/command/aggregate"
	"github.com/payfazz/go-apt/pkg/esfazz"
)

// AccountEventRepository is repository for account event
type AccountEventRepository interface {
	Save(ctx context.Context, payload *esfazz.EventPayload) (*aggregate.Account, error)
	Find(ctx context.Context, id string) (*aggregate.Account, error)
}

type accountEventRepository struct {
	repository esfazz.EventRepository
}

// Save save account event and aggregate snapshot to storage
func (a *accountEventRepository) Save(ctx context.Context, payload *esfazz.EventPayload) (*aggregate.Account, error) {
	// save to event store
	result, err := a.repository.Save(ctx, payload)
	if err != nil {
		return nil, err
	}
	return result.(*aggregate.Account), nil
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
func NewAccountEventRepository() AccountEventRepository {
	return &accountEventRepository{
		repository: esfazz.NewEventRepository(&aggregate.Account{}, "account_event", "account_aggregate"),
	}
}
