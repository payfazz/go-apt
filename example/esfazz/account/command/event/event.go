package event

import (
	"github.com/payfazz/go-apt/pkg/esfazz"
	"time"
)

// AccountCreated create account created event
func AccountCreated(id string, name string, balance int64) *esfazz.EventPayload {
	return &esfazz.EventPayload{
		Type:      AccountCreatedType,
		Aggregate: NewAccount(id),
		Data: AccountCreatedData{
			Name:      name,
			Balance:   balance,
			CreatedAt: time.Now(),
		},
	}
}

// AccountNameChanged create account name changed event
func AccountNameChanged(account *Account, name string) *esfazz.EventPayload {
	return &esfazz.EventPayload{
		Type:      AccountNameChangedType,
		Aggregate: account,
		Data: AccountNameChangedData{
			Name:      name,
			UpdatedAt: time.Now(),
		},
	}
}

// AccountDeposited create account deposited event
func AccountDeposited(account *Account, amount int64) *esfazz.EventPayload {
	return &esfazz.EventPayload{
		Type:      AccountDepositedType,
		Aggregate: account,
		Data: AccountDepositedData{
			Amount:    amount,
			UpdatedAt: time.Now(),
		},
	}
}

// AccountWithdrawn create account deposited event
func AccountWithdrawn(account *Account, amount int64) *esfazz.EventPayload {
	return &esfazz.EventPayload{
		Type:      AccountWithdrawnType,
		Aggregate: account,
		Data: AccountWithdrawnData{
			Amount:    amount,
			UpdatedAt: time.Now(),
		},
	}
}

// AccountDeleted create account deposited event
func AccountDeleted(account *Account) *esfazz.EventPayload {
	return &esfazz.EventPayload{
		Type:      AccountDeletedType,
		Aggregate: account,
		Data: AccountDeletedData{
			DeletedAt: time.Now(),
		},
	}
}
