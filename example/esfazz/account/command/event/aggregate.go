package event

import (
	"encoding/json"
	"github.com/payfazz/go-apt/pkg/esfazz"
	"time"
)

// Account is aggregate object for account
type Account struct {
	esfazz.BaseAggregate
	Name      string     `json:"name"`
	Balance   int64      `json:"balance"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// Apply apply event to aggregate
func (a *Account) Apply(evs ...*esfazz.Event) error {
	for _, ev := range evs {
		a.Version = a.Version + 1
		var err error
		switch ev.Type {
		case AccountCreatedType:
			err = a.applyAccountCreated(ev)
		case AccountNameChangedType:
			err = a.applyAccountNameChanged(ev)
		case AccountDepositedType:
			err = a.applyAccountDeposited(ev)
		case AccountWithdrawnType:
			err = a.applyAccountWithdrawn(ev)
		case AccountDeletedType:
			err = a.applyAccountDeleted(ev)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Account) applyAccountCreated(ev *esfazz.Event) error {
	data := &AccountCreatedData{}
	err := json.Unmarshal(ev.Data, data)
	if err != nil {
		return err
	}

	a.Id = ev.Aggregate.GetId()
	a.Name = data.Name
	a.Balance = data.Balance
	return nil
}

func (a *Account) applyAccountNameChanged(ev *esfazz.Event) error {
	data := &AccountNameChangedData{}
	err := json.Unmarshal(ev.Data, data)
	if err != nil {
		return err
	}

	a.Name = data.Name
	return nil
}

func (a *Account) applyAccountDeposited(ev *esfazz.Event) error {
	data := &AccountDepositedData{}
	err := json.Unmarshal(ev.Data, data)
	if err != nil {
		return err
	}

	a.Balance += data.Amount
	return nil
}

func (a *Account) applyAccountWithdrawn(ev *esfazz.Event) error {
	data := &AccountWithdrawnData{}
	err := json.Unmarshal(ev.Data, data)
	if err != nil {
		return err
	}

	a.Balance -= data.Amount
	return nil

}

func (a *Account) applyAccountDeleted(ev *esfazz.Event) error {
	data := &AccountDeletedData{}
	err := json.Unmarshal(ev.Data, data)
	if err != nil {
		return err
	}

	a.DeletedAt = &data.DeletedAt
	return nil
}

// NewAccount is constructor for account aggregate
func NewAccount(id string) esfazz.Aggregate {
	acc := &Account{}
	acc.Id = id
	return acc
}
