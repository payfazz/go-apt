package repository

import (
	"context"
	"encoding/json"
	"github.com/payfazz/go-apt/example/esfazz/domain/account/command/aggregate"
	"github.com/payfazz/go-apt/pkg/esfazz"
)

type AccountEventRepository interface {
	Save(ctx context.Context, payload esfazz.EventPayload) (*aggregate.Account, error)
	Find(ctx context.Context, id string) (*aggregate.Account, error)
}

type accountEventRepository struct {
	eventStore esfazz.EventStore
	aggStore   esfazz.AggregateStore
}

func (a *accountEventRepository) Save(ctx context.Context, payload esfazz.EventPayload) (*aggregate.Account, error) {
	// save to event store
	savedEvent, err := a.eventStore.Save(ctx, payload)
	if err != nil {
		return nil, err
	}

	account, err := a.saveAggregate(ctx, savedEvent.AggregateId)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (a *accountEventRepository) Find(ctx context.Context, id string) (*aggregate.Account, error) {
	account := &aggregate.Account{}

	// load data from saved aggregate
	agg, err := a.aggStore.FindBy(ctx, id)
	if err != nil {
		return nil, err
	}
	if agg != nil {
		err = json.Unmarshal(agg.Data, account)
		if err != nil {
			return nil, err
		}
	}

	// load new event and apply
	evs, err := a.eventStore.FindAllBy(ctx, id, account.Version)
	if err != nil {
		return nil, err
	}
	err = account.ApplyAll(evs...)
	if err != nil {
		return nil, err
	}

	// return nil if no event applied
	if account.Version == 0 {
		return nil, nil
	}

	return account, nil
}

func (a *accountEventRepository) saveAggregate(ctx context.Context, id string) (*aggregate.Account, error) {
	// save snapshot
	account, err := a.Find(ctx, id)
	if err != nil {
		return nil, err
	}
	_, err = a.aggStore.Save(ctx, account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func NewAccountEventRepository() AccountEventRepository {
	return &accountEventRepository{
		eventStore: esfazz.PostgresEventStore("account_event"),
		aggStore:   esfazz.PostgresAggregateStore("account_aggregate"),
	}
}
