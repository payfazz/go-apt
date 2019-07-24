package repository

import (
	"context"
	"encoding/json"
	"github.com/payfazz/go-apt/example/eventsource/domain/account/command/aggregate"
	"github.com/payfazz/go-apt/pkg/esfazz"
)

type AccountEventRepository interface {
	Save(ctx context.Context, payload esfazz.EventPayload) (*aggregate.Account, error)
	Find(ctx context.Context, id string) (*aggregate.Account, error)
}

type accountEventRepository struct {
	eventStore esfazz.EventStore
	snapStore  esfazz.SnapshotStore
}

func (a *accountEventRepository) Save(ctx context.Context, payload esfazz.EventPayload) (*aggregate.Account, error) {
	// save to event store
	savedEvent, err := a.eventStore.Save(ctx, payload)
	if err != nil {
		return nil, err
	}

	account, err := a.saveSnapshot(ctx, savedEvent.AggregateId)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (a *accountEventRepository) Find(ctx context.Context, id string) (*aggregate.Account, error) {
	account := &aggregate.Account{}

	// load data from snapshot
	snap, err := a.snapStore.FindBy(ctx, id)
	if err != nil {
		return nil, err
	}
	if snap != nil {
		err = json.Unmarshal(snap.Data, account)
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

func (a *accountEventRepository) saveSnapshot(ctx context.Context, id string) (*aggregate.Account, error) {
	// save snapshot
	account, err := a.Find(ctx, id)
	if err != nil {
		return nil, err
	}
	_, err = a.snapStore.Save(ctx, account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func NewAccountEventRepository() AccountEventRepository {
	return &accountEventRepository{
		eventStore: esfazz.PostgresEventStore("account_event"),
		snapStore:  esfazz.PostgresSnapshotStore("account_snapshot"),
	}
}
