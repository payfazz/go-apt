package command

import (
	"context"
	"errors"
	"github.com/payfazz/go-apt/example/esfazz/domain/account/command/aggregate"
	"github.com/payfazz/go-apt/example/esfazz/domain/account/command/data"
	"github.com/payfazz/go-apt/example/esfazz/domain/account/command/event"
	"github.com/payfazz/go-apt/example/esfazz/domain/account/command/repository"
	"github.com/payfazz/go-apt/pkg/esfazz"
)

// AccountCommand is an interface for account command
type AccountCommand interface {
	Create(ctx context.Context, payload data.CreatePayload) (*aggregate.Account, error)
	ChangeName(ctx context.Context, payload data.ChangeNamePayload) (*aggregate.Account, error)
	Deposit(ctx context.Context, payload data.DepositPayload) (*aggregate.Account, error)
	Withdraw(ctx context.Context, payload data.WithdrawPayload) (*aggregate.Account, error)
	Delete(ctx context.Context, accountId string) (*aggregate.Account, error)
}

type accountCommand struct {
	repository repository.AccountEventRepository
}

// Create is command to create account
func (a *accountCommand) Create(ctx context.Context, payload data.CreatePayload) (*aggregate.Account, error) {
	ev := &esfazz.EventPayload{
		Type: event.ACCOUNT_CREATED,
		Data: event.AccountCreatedData{
			Name:    payload.Name,
			Balance: payload.Balance,
		},
	}

	account, err := a.repository.Save(ctx, ev)
	if err != nil {
		return nil, err
	}

	return account, nil
}

// ChangeName is command to change account name
func (a *accountCommand) ChangeName(ctx context.Context, payload data.ChangeNamePayload) (*aggregate.Account, error) {
	account, err := a.repository.Find(ctx, payload.AccountId)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return nil, errors.New("account not found")
	}
	if account.DeletedAt != nil {
		return nil, errors.New("account deleted")
	}

	ev := &esfazz.EventPayload{
		Aggregate: account,
		Type:      event.ACCOUNT_NAME_CHANGED,
		Data: event.AccountNameChangedData{
			Name: payload.Name,
		},
	}

	account, err = a.repository.Save(ctx, ev)
	if err != nil {
		return nil, err
	}
	return account, nil
}

// Deposit is command to create account deposit
func (a *accountCommand) Deposit(ctx context.Context, payload data.DepositPayload) (*aggregate.Account, error) {
	account, err := a.repository.Find(ctx, payload.AccountId)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return nil, errors.New("account not found")
	}
	if account.DeletedAt != nil {
		return nil, errors.New("account deleted")
	}

	ev := &esfazz.EventPayload{
		Aggregate: account,
		Type:      event.ACCOUNT_DEPOSITED,
		Data: event.AccountDepositedData{
			Amount: payload.Amount,
		},
	}

	account, err = a.repository.Save(ctx, ev)
	if err != nil {
		return nil, err
	}
	return account, nil
}

// Withdraw is command to create account withdraw
func (a *accountCommand) Withdraw(ctx context.Context, payload data.WithdrawPayload) (*aggregate.Account, error) {
	account, err := a.repository.Find(ctx, payload.AccountId)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return nil, errors.New("account not found")
	}
	if account.DeletedAt != nil {
		return nil, errors.New("account deleted")
	}
	if account.Balance < payload.Amount {
		return nil, errors.New("account balance is smaller than withdraw ammount")
	}

	ev := &esfazz.EventPayload{
		Aggregate: account,
		Type:      event.ACCOUNT_WITHDRAWN,
		Data: event.AccountWithdrawnData{
			Amount: payload.Amount,
		},
	}

	account, err = a.repository.Save(ctx, ev)
	if err != nil {
		return nil, err
	}
	return account, nil
}

// Delete is command to delete account
func (a *accountCommand) Delete(ctx context.Context, accountId string) (*aggregate.Account, error) {
	account, err := a.repository.Find(ctx, accountId)
	if err != nil {
		return nil, err
	}

	if account == nil {
		return nil, errors.New("account not found")
	}
	if account.DeletedAt != nil {
		// account already deleted
		return nil, nil
	}
	if account.Balance != 0 {
		return nil, errors.New("account balance must be zero before deleted")
	}

	ev := &esfazz.EventPayload{
		Aggregate: account,
		Type:      event.ACCOUNT_DELETED,
	}

	account, err = a.repository.Save(ctx, ev)
	if err != nil {
		return nil, err
	}
	return account, nil
}

// NewAccountCommand create new account command service
func NewAccountCommand() AccountCommand {
	return &accountCommand{
		repository: repository.NewAccountEventRepository(),
	}
}
