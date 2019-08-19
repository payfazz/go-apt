package model

import (
	"github.com/payfazz/go-apt/pkg/esfazz"
	"time"
)

// Created create account created event
func (a *Account) Created(name string, balance int64) *esfazz.EventPayload {
	return &esfazz.EventPayload{
		Type:      AccountCreatedType,
		Aggregate: a,
		Data: AccountCreatedData{
			Name:      name,
			Balance:   balance,
			CreatedAt: time.Now(),
		},
	}
}

// NameChanged create account name changed event
func (a *Account) NameChanged(name string) *esfazz.EventPayload {
	return &esfazz.EventPayload{
		Type:      AccountNameChangedType,
		Aggregate: a,
		Data: AccountNameChangedData{
			Name:      name,
			UpdatedAt: time.Now(),
		},
	}
}

// Deposited create account deposited event
func (a *Account) Deposited(amount int64) *esfazz.EventPayload {
	return &esfazz.EventPayload{
		Type:      AccountDepositedType,
		Aggregate: a,
		Data: AccountDepositedData{
			Amount:    amount,
			UpdatedAt: time.Now(),
		},
	}
}

// Withdrawn create account deposited event
func (a *Account) Withdrawn(amount int64) *esfazz.EventPayload {
	return &esfazz.EventPayload{
		Type:      AccountWithdrawnType,
		Aggregate: a,
		Data: AccountWithdrawnData{
			Amount:    amount,
			UpdatedAt: time.Now(),
		},
	}
}

// Deleted create account deposited event
func (a *Account) Deleted() *esfazz.EventPayload {
	return &esfazz.EventPayload{
		Type:      AccountDeletedType,
		Aggregate: a,
		Data: AccountDeletedData{
			DeletedAt: time.Now(),
		},
	}
}
