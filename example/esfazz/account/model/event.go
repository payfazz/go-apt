package model

import (
	"github.com/payfazz/go-apt/pkg/esfazz"
	"time"
)

// AccountCreated create account created event
func AccountCreated(name string, balance int64) *esfazz.EventPayload {
	return &esfazz.EventPayload{
		Type: AccountCreatedType,
		Data: AccountCreatedData{
			Name:      name,
			Balance:   balance,
			CreatedAt: time.Now(),
		},
	}
}

// AccountNameChanged create account name changed event
func AccountNameChanged(name string) *esfazz.EventPayload {
	return &esfazz.EventPayload{
		Type: AccountNameChangedType,
		Data: AccountNameChangedData{
			Name:      name,
			UpdatedAt: time.Now(),
		},
	}
}

// AccountDeposited create account deposited event
func AccountDeposited(amount int64) *esfazz.EventPayload {
	return &esfazz.EventPayload{
		Type: AccountDepositedType,
		Data: AccountDepositedData{
			Amount:    amount,
			UpdatedAt: time.Now(),
		},
	}
}

// AccountWithdrawn create account deposited event
func AccountWithdrawn(amount int64) *esfazz.EventPayload {
	return &esfazz.EventPayload{
		Type: AccountWithdrawnType,
		Data: AccountWithdrawnData{
			Amount:    amount,
			UpdatedAt: time.Now(),
		},
	}
}

// AccountDeleted create account deposited event
func AccountDeleted() *esfazz.EventPayload {
	return &esfazz.EventPayload{
		Type: AccountDeletedType,
		Data: AccountDeletedData{
			DeletedAt: time.Now(),
		},
	}
}
