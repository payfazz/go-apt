package event

import (
	"github.com/gofrs/uuid"
	"github.com/payfazz/go-apt/pkg/esfazz"
	"time"
)

// AccountCreated create account created event
func AccountCreated(name string, balance int64) (*esfazz.Event, error) {
	uuidV4, _ := uuid.NewV4()

	payload := &esfazz.EventPayload{
		Type: AccountCreatedType,
		Aggregate: &esfazz.BaseAggregate{
			Id:      uuidV4.String(),
			Version: 0,
		},
		Data: AccountCreatedData{
			Name:      name,
			Balance:   balance,
			CreatedAt: time.Now(),
		},
	}
	return esfazz.CreateEvent(payload)
}

// AccountNameChanged create account name changed event
func AccountNameChanged(agg esfazz.Aggregate, name string) (*esfazz.Event, error) {
	payload := &esfazz.EventPayload{
		Type:      AccountNameChangedType,
		Aggregate: agg,
		Data: AccountNameChangedData{
			Name:      name,
			UpdatedAt: time.Now(),
		},
	}
	return esfazz.CreateEvent(payload)

}

// AccountDeposited create account deposited event
func AccountDeposited(agg esfazz.Aggregate, amount int64) (*esfazz.Event, error) {
	payload := &esfazz.EventPayload{
		Type:      AccountDepositedType,
		Aggregate: agg,
		Data: AccountDepositedData{
			Amount:    amount,
			UpdatedAt: time.Now(),
		},
	}
	return esfazz.CreateEvent(payload)
}

// AccountWithdrawn create account deposited event
func AccountWithdrawn(agg esfazz.Aggregate, amount int64) (*esfazz.Event, error) {
	payload := &esfazz.EventPayload{
		Type:      AccountWithdrawnType,
		Aggregate: agg,
		Data: AccountWithdrawnData{
			Amount:    amount,
			UpdatedAt: time.Now(),
		},
	}
	return esfazz.CreateEvent(payload)
}

// AccountDeleted create account deposited event
func AccountDeleted(agg esfazz.Aggregate) (*esfazz.Event, error) {
	payload := &esfazz.EventPayload{
		Type:      AccountDeletedType,
		Aggregate: agg,
		Data: AccountDeletedData{
			DeletedAt: time.Now(),
		},
	}
	return esfazz.CreateEvent(payload)
}

const (
	// AccountCreatedType is string type for account created
	AccountCreatedType = "account.created"
	// AccountNameChangedType is string type for account name changed
	AccountNameChangedType = "account.name.changed"
	// AccountDepositedType is string type for account deposited
	AccountDepositedType = "account.deposited"
	// AccountWithdrawnType is string type for account withdrawn
	AccountWithdrawnType = "account.withdrawn"
	// AccountDeletedType is string type for account deleted
	AccountDeletedType = "account.deleted"
)

// AccountCreatedData is account event created data
type AccountCreatedData struct {
	Name      string    `json:"name"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

// AccountNameChangedData is account name changed event data
type AccountNameChangedData struct {
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AccountDepositedData is account deposited event data
type AccountDepositedData struct {
	Amount    int64     `json:"amount"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AccountWithdrawnData is account withdrawn event data
type AccountWithdrawnData struct {
	Amount    int64     `json:"amount"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AccountDeletedData is account deleted event data
type AccountDeletedData struct {
	DeletedAt time.Time `json:"deleted_at"`
}
