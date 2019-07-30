package aggregate

import (
	"encoding/json"
	"github.com/payfazz/go-apt/example/esfazz/domain/account/command/event"
	"github.com/payfazz/go-apt/pkg/esfazz"
	"time"
)

// Account is aggregate object for account
type Account struct {
	esfazz.BaseAggregate
	Name      string     `json:"name"`
	Balance   int64      `json:"balance"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// Apply apply event to aggregate
func (a *Account) Apply(eventLog *esfazz.EventLog) error {
	a.Version = a.Version + 1
	a.UpdatedAt = eventLog.CreatedAt

	switch eventLog.EventType {
	case event.AccountCreatedType:
		return a.applyAccountCreated(eventLog)
	case event.AccountNameChangedType:
		return a.applyAccountNameChanged(eventLog)
	case event.AccountDepositedType:
		return a.applyAccountDeposited(eventLog)
	case event.AccountWithdrawnType:
		return a.applyAccountWithdrawn(eventLog)
	case event.AccountDeletedType:
		return a.applyAccountDeleted(eventLog)
	}

	return nil
}

func (a *Account) applyAccountCreated(eventLog *esfazz.EventLog) error {
	data := &event.AccountCreatedData{}
	err := json.Unmarshal(eventLog.Data, data)
	if err != nil {
		return err
	}

	a.Id = eventLog.AggregateId
	a.CreatedAt = eventLog.CreatedAt
	a.Name = data.Name
	a.Balance = data.Balance
	return nil
}

func (a *Account) applyAccountNameChanged(eventLog *esfazz.EventLog) error {
	data := &event.AccountNameChangedData{}
	err := json.Unmarshal(eventLog.Data, data)
	if err != nil {
		return err
	}

	a.Name = data.Name
	return nil
}

func (a *Account) applyAccountDeposited(eventLog *esfazz.EventLog) error {
	data := &event.AccountDepositedData{}
	err := json.Unmarshal(eventLog.Data, data)
	if err != nil {
		return err
	}

	a.Balance += data.Amount
	return nil
}

func (a *Account) applyAccountWithdrawn(eventLog *esfazz.EventLog) error {
	data := &event.AccountWithdrawnData{}
	err := json.Unmarshal(eventLog.Data, data)
	if err != nil {
		return err
	}

	a.Balance -= data.Amount
	return nil

}

func (a *Account) applyAccountDeleted(eventLog *esfazz.EventLog) error {
	a.DeletedAt = eventLog.CreatedAt
	return nil
}
