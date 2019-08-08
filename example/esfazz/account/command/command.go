package command

import (
	"context"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/payfazz/go-apt/example/esfazz/account/command/event"
	"github.com/payfazz/go-apt/pkg/esfazz/esrepo"
)

// AccountCommand is an interface for account command
type AccountCommand interface {
	Create(ctx context.Context, payload CreatePayload) (*event.Account, error)
	ChangeName(ctx context.Context, payload ChangeNamePayload) (*event.Account, error)
	Deposit(ctx context.Context, payload DepositPayload) (*event.Account, error)
	Withdraw(ctx context.Context, payload WithdrawPayload) (*event.Account, error)
	Delete(ctx context.Context, accountId string) (*event.Account, error)
}

type accountCommand struct {
	repository esrepo.Repository
}

// Create is command to create account
func (a *accountCommand) Create(ctx context.Context, payload CreatePayload) (*event.Account, error) {
	uuidV4, _ := uuid.NewV4()

	ev := event.AccountCreated(uuidV4.String(), payload.Name, payload.Balance)
	err := a.repository.Save(ctx, ev)
	if err != nil {
		return nil, err
	}

	agg, err := a.repository.Find(ctx, ev.Aggregate.GetId())
	account := agg.(*event.Account)
	return account, nil
}

// ChangeName is command to change account name
func (a *accountCommand) ChangeName(ctx context.Context, payload ChangeNamePayload) (*event.Account, error) {
	agg, err := a.repository.Find(ctx, payload.AccountId)
	if err != nil {
		return nil, err
	}

	account := agg.(*event.Account)
	if account == nil || account.DeletedAt != nil {
		return nil, errors.New("account not found or deleted")
	}

	ev := event.AccountNameChanged(account, payload.Name)
	err = a.repository.Save(ctx, ev)
	if err != nil {
		return nil, err
	}

	agg, err = a.repository.Find(ctx, ev.Aggregate.GetId())
	if err != nil {
		return nil, err
	}
	account = agg.(*event.Account)
	return account, nil
}

// Deposit is command to create account deposit
func (a *accountCommand) Deposit(ctx context.Context, payload DepositPayload) (*event.Account, error) {
	agg, err := a.repository.Find(ctx, payload.AccountId)
	if err != nil {
		return nil, err
	}

	account := agg.(*event.Account)
	if account == nil || account.DeletedAt != nil {
		return nil, errors.New("account not found or deleted")
	}

	ev := event.AccountDeposited(account, payload.Amount)
	err = a.repository.Save(ctx, ev)
	if err != nil {
		return nil, err
	}

	agg, err = a.repository.Find(ctx, ev.Aggregate.GetId())
	if err != nil {
		return nil, err
	}
	account = agg.(*event.Account)
	return account, nil
}

// Withdraw is command to create account withdraw
func (a *accountCommand) Withdraw(ctx context.Context, payload WithdrawPayload) (*event.Account, error) {
	agg, err := a.repository.Find(ctx, payload.AccountId)
	if err != nil {
		return nil, err
	}
	account := agg.(*event.Account)
	if account == nil || account.DeletedAt != nil {
		return nil, errors.New("account not found or deleted")
	}
	if account.Balance < payload.Amount {
		return nil, errors.New("account balance is smaller than withdraw amount")
	}

	ev := event.AccountWithdrawn(account, payload.Amount)
	err = a.repository.Save(ctx, ev)
	if err != nil {
		return nil, err
	}

	agg, err = a.repository.Find(ctx, ev.Aggregate.GetId())
	if err != nil {
		return nil, err
	}
	account = agg.(*event.Account)
	return account, nil
}

// Delete is command to delete account
func (a *accountCommand) Delete(ctx context.Context, accountId string) (*event.Account, error) {
	agg, err := a.repository.Find(ctx, accountId)
	if err != nil {
		return nil, err
	}

	account := agg.(*event.Account)
	if account == nil || account.DeletedAt != nil {
		return nil, errors.New("account not found")
	}
	if account.Balance != 0 {
		return nil, errors.New("account balance must be zero before deleted")
	}

	ev := event.AccountDeleted(account)
	err = a.repository.Save(ctx, ev)
	if err != nil {
		return nil, err
	}

	agg, err = a.repository.Find(ctx, ev.Aggregate.GetId())
	if err != nil {
		return nil, err
	}
	account = agg.(*event.Account)
	return account, nil
}

// NewAccountCommand create new account command service
func NewAccountCommand(repository esrepo.Repository) AccountCommand {
	return &accountCommand{
		repository: repository,
	}
}
