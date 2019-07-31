package aggregate

import (
	"encoding/json"
	"github.com/payfazz/go-apt/example/esfazz/domain/account/command/event"
	"github.com/payfazz/go-apt/pkg/esfazz"
	"time"
)

// Account is aggregate object for account
type Account struct {
	Id        string     `json:"id"`
	Version   int        `json:"version"`
	Name      string     `json:"name"`
	Balance   int64      `json:"balance"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// GetId return id of account
func (a *Account) GetId() string {
	return a.Id
}

// GetVersion return version of account
func (a *Account) GetVersion() int {
	return a.Version
}

// Apply apply event to aggregate
func (a *Account) Apply(ev *esfazz.Event) error {
	a.Version = a.Version + 1

	switch ev.Type {
	case event.AccountCreatedType:
		return a.applyAccountCreated(ev)
	case event.AccountNameChangedType:
		return a.applyAccountNameChanged(ev)
	case event.AccountDepositedType:
		return a.applyAccountDeposited(ev)
	case event.AccountWithdrawnType:
		return a.applyAccountWithdrawn(ev)
	case event.AccountDeletedType:
		return a.applyAccountDeleted(ev)
	}

	return nil
}

func (a *Account) applyAccountCreated(ev *esfazz.Event) error {
	data := &event.AccountCreatedData{}
	err := json.Unmarshal(ev.Data, data)
	if err != nil {
		return err
	}

	a.Id = ev.Aggregate.GetId()
	a.CreatedAt = &data.CreatedAt
	a.UpdatedAt = a.CreatedAt
	a.Name = data.Name
	a.Balance = data.Balance
	return nil
}

func (a *Account) applyAccountNameChanged(ev *esfazz.Event) error {
	data := &event.AccountNameChangedData{}
	err := json.Unmarshal(ev.Data, data)
	if err != nil {
		return err
	}

	a.Name = data.Name
	a.UpdatedAt = &data.UpdatedAt
	return nil
}

func (a *Account) applyAccountDeposited(ev *esfazz.Event) error {
	data := &event.AccountDepositedData{}
	err := json.Unmarshal(ev.Data, data)
	if err != nil {
		return err
	}

	a.Balance += data.Amount
	a.UpdatedAt = &data.UpdatedAt
	return nil
}

func (a *Account) applyAccountWithdrawn(ev *esfazz.Event) error {
	data := &event.AccountWithdrawnData{}
	err := json.Unmarshal(ev.Data, data)
	if err != nil {
		return err
	}

	a.Balance -= data.Amount
	a.UpdatedAt = &data.UpdatedAt
	return nil

}

func (a *Account) applyAccountDeleted(ev *esfazz.Event) error {
	data := &event.AccountDeletedData{}
	err := json.Unmarshal(ev.Data, data)
	if err != nil {
		return err
	}

	a.DeletedAt = &data.DeletedAt
	return nil
}

// AccountAggregate is constructor for account aggregate
func AccountAggregate(id string) esfazz.Aggregate {
	return &Account{
		Id: id,
	}
}
