package command

import (
	"context"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/payfazz/go-apt/example/esfazz/account/model"
	"github.com/payfazz/go-apt/pkg/esfazz/esrepo"
)

// AccountCommand is an interface for account command
type AccountCommand interface {
	Create(ctx context.Context, payload CreatePayload) (*string, error)
	ChangeName(ctx context.Context, payload ChangeNamePayload) error
	Deposit(ctx context.Context, payload DepositPayload) error
	Withdraw(ctx context.Context, payload WithdrawPayload) error
	Delete(ctx context.Context, accountId string) error
}

type accountCommand struct {
	repository esrepo.Repository
}

// Create is command to create account
func (a *accountCommand) Create(ctx context.Context, payload CreatePayload) (*string, error) {
	uuidV4, _ := uuid.NewV4()
	id := uuidV4.String()
	account := model.NewAccount(id).(*model.Account)

	err := a.repository.Save(ctx, account, model.AccountCreated(payload.Name, payload.Balance))
	if err != nil {
		return nil, err
	}

	return &id, nil
}

// ChangeName is command to change account name
func (a *accountCommand) ChangeName(ctx context.Context, payload ChangeNamePayload) error {
	agg, err := a.repository.Find(ctx, payload.AccountId)
	if err != nil {
		return err
	}

	account := agg.(*model.Account)
	if account == nil || account.DeletedTime != nil {
		return errors.New("account not found or deleted")
	}

	err = a.repository.Save(ctx, account, model.AccountNameChanged(payload.Name))
	return err
}

// Deposit is command to create account deposit
func (a *accountCommand) Deposit(ctx context.Context, payload DepositPayload) error {
	agg, err := a.repository.Find(ctx, payload.AccountId)
	if err != nil {
		return err
	}

	account := agg.(*model.Account)
	if account == nil || account.DeletedTime != nil {
		return errors.New("account not found or deleted")
	}

	err = a.repository.Save(ctx, account, model.AccountDeposited(payload.Amount))
	return err
}

// Withdraw is command to create account withdraw
func (a *accountCommand) Withdraw(ctx context.Context, payload WithdrawPayload) error {
	agg, err := a.repository.Find(ctx, payload.AccountId)
	if err != nil {
		return err
	}
	account := agg.(*model.Account)
	if account == nil || account.DeletedTime != nil {
		return errors.New("account not found or deleted")
	}
	if account.Balance < payload.Amount {
		return errors.New("account balance is smaller than withdraw amount")
	}

	err = a.repository.Save(ctx, account, model.AccountWithdrawn(payload.Amount))
	return err

}

// Delete is command to delete account
func (a *accountCommand) Delete(ctx context.Context, accountId string) error {
	agg, err := a.repository.Find(ctx, accountId)
	if err != nil {
		return err
	}

	account := agg.(*model.Account)
	if account == nil || account.DeletedTime != nil {
		return errors.New("account not found")
	}
	if account.Balance != 0 {
		return errors.New("account balance must be zero before deleted")
	}

	err = a.repository.Save(ctx, account, model.AccountDeleted())
	return err
}

// NewAccountCommand create new account command service
func NewAccountCommand(repository esrepo.Repository) AccountCommand {
	return &accountCommand{
		repository: repository,
	}
}
