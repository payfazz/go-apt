package model

import (
	"encoding/json"
	"github.com/payfazz/go-apt/pkg/esfazz"
	"github.com/payfazz/go-apt/pkg/fazzdb"
	"time"
)

// Account is aggregate object for account
type Account struct {
	fazzdb.Model `json:"-"`
	esfazz.BaseAggregate
	Name        string     `json:"name" db:"name"`
	Balance     int64      `json:"balance" db:"balance"`
	DeletedTime *time.Time `json:"deleted_time" db:"deleted_time"`
}

// Get used to get the data
func (a *Account) Get(key string) interface{} { return a.Payload()[key] }

// Payload used to get the payload
func (a *Account) Payload() map[string]interface{} {
	results := a.MapPayload(a)
	results["id"] = a.BaseAggregate.Id
	results["version"] = a.BaseAggregate.Version
	return results
}

// Apply apply event to aggregate
func (a *Account) Apply(evs ...*esfazz.Event) error {
	for _, ev := range evs {
		var err error
		switch ev.Type {
		case AccountCreatedType:
			a.Version = a.Version + 1
			err = a.applyCreated(ev)
		case AccountNameChangedType:
			a.Version = a.Version + 1
			err = a.applyNameChanged(ev)
		case AccountDepositedType:
			a.Version = a.Version + 1
			err = a.applyDeposited(ev)
		case AccountWithdrawnType:
			a.Version = a.Version + 1
			err = a.applyWithdrawn(ev)
		case AccountDeletedType:
			a.Version = a.Version + 1
			err = a.applyDeleted(ev)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Account) applyCreated(ev *esfazz.Event) error {
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

func (a *Account) applyNameChanged(ev *esfazz.Event) error {
	data := &AccountNameChangedData{}
	err := json.Unmarshal(ev.Data, data)
	if err != nil {
		return err
	}

	a.Name = data.Name
	return nil
}

func (a *Account) applyDeposited(ev *esfazz.Event) error {
	data := &AccountDepositedData{}
	err := json.Unmarshal(ev.Data, data)
	if err != nil {
		return err
	}

	a.Balance += data.Amount
	return nil
}

func (a *Account) applyWithdrawn(ev *esfazz.Event) error {
	data := &AccountWithdrawnData{}
	err := json.Unmarshal(ev.Data, data)
	if err != nil {
		return err
	}

	a.Balance -= data.Amount
	return nil

}

func (a *Account) applyDeleted(ev *esfazz.Event) error {
	data := &AccountDeletedData{}
	err := json.Unmarshal(ev.Data, data)
	if err != nil {
		return err
	}

	a.DeletedTime = &data.DeletedAt
	return nil
}

// NewAccount is constructor for account aggregate
func NewAccount(id string) esfazz.Aggregate {
	acc := AccountModel()
	acc.Id = id
	return acc
}

// AccountModel is a constructor for account model
func AccountModel() *Account {
	return &Account{
		Model: fazzdb.PlainModel(
			"account",
			[]fazzdb.Column{
				fazzdb.Col("id"),
				fazzdb.Col("version"),
				fazzdb.Col("name"),
				fazzdb.Col("balance"),
				fazzdb.Col("deleted_time"),
			},
			"id",
			false,
			false,
		),
	}
}
