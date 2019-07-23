package command

import (
	"context"
	"errors"
	"github.com/payfazz/go-apt/example/eventsource/domain/account/command/aggregate"
	"github.com/payfazz/go-apt/example/eventsource/domain/account/command/data"
	"github.com/payfazz/go-apt/example/eventsource/domain/account/command/event"
	"github.com/payfazz/go-apt/example/eventsource/domain/account/command/repository"
	"github.com/payfazz/go-apt/pkg/esfazz"
)

type AccountCommand interface {
	Create(ctx context.Context, payload data.CreatePayload) (*string, error)
	ChangeName(ctx context.Context, payload data.ChangeNamePayload) error
	Deposit(ctx context.Context, payload data.DepositPayload) error
	Withdraw(ctx context.Context, payload data.WithdrawPayload) error
	Delete(ctx context.Context, accountId string) error
	DirectGet(ctx context.Context, accountId string) (*aggregate.Account, error)
}

type accountCommand struct {
	repository repository.AccountEventRepository
}

func (a *accountCommand) Create(ctx context.Context, payload data.CreatePayload) (*string, error) {
	ev := esfazz.EventPayload{
		Type: event.ACCOUNT_CREATED,
		Data: event.AccountCreatedData{
			Name:    payload.Name,
			Balance: payload.Balance,
		},
	}

	evLog, err := a.repository.Save(ctx, ev)
	if err != nil {
		return nil, err
	}

	return &evLog.AggregateId, nil
}

func (a *accountCommand) ChangeName(ctx context.Context, payload data.ChangeNamePayload) error {
	account, err := a.repository.Find(ctx, payload.AccountId)
	if err != nil {
		return err
	}

	if account == nil {
		return errors.New("account not found")
	}
	if account.DeletedAt != nil {
		return errors.New("account deleted")
	}

	ev := esfazz.EventPayload{
		Aggregate: account,
		Type:      event.ACCOUNT_NAME_CHANGED,
		Data: event.AccountNameChangedData{
			Name: payload.Name,
		},
	}

	_, err = a.repository.Save(ctx, ev)
	return err
}

func (a *accountCommand) Deposit(ctx context.Context, payload data.DepositPayload) error {
	account, err := a.repository.Find(ctx, payload.AccountId)
	if err != nil {
		return err
	}

	if account == nil {
		return errors.New("account not found")
	}
	if account.DeletedAt != nil {
		return errors.New("account deleted")
	}

	ev := esfazz.EventPayload{
		Aggregate: account,
		Type:      event.ACCOUNT_DEPOSITED,
		Data: event.AccountDepositedData{
			Amount: payload.Amount,
		},
	}

	_, err = a.repository.Save(ctx, ev)
	return err
}

func (a *accountCommand) Withdraw(ctx context.Context, payload data.WithdrawPayload) error {
	account, err := a.repository.Find(ctx, payload.AccountId)
	if err != nil {
		return err
	}

	if account == nil {
		return errors.New("account not found")
	}
	if account.DeletedAt != nil {
		return errors.New("account deleted")
	}
	if account.Balance < payload.Amount {
		return errors.New("account balance is smaller than withdraw ammount")
	}

	ev := esfazz.EventPayload{
		Aggregate: account,
		Type:      event.ACCOUNT_WITHDRAWN,
		Data: event.AccountWithdrawnData{
			Amount: payload.Amount,
		},
	}

	_, err = a.repository.Save(ctx, ev)
	return err
}

func (a *accountCommand) Delete(ctx context.Context, accountId string) error {
	account, err := a.repository.Find(ctx, accountId)
	if err != nil {
		return err
	}

	if account == nil {
		return errors.New("account not found")
	}
	if account.DeletedAt != nil {
		// account already deleted
		return nil
	}
	if account.Balance != 0 {
		return errors.New("account balance must be zero before deleted")
	}

	ev := esfazz.EventPayload{
		Aggregate: account,
		Type:      event.ACCOUNT_DELETED,
	}

	_, err = a.repository.Save(ctx, ev)
	return err
}

func (a *accountCommand) DirectGet(ctx context.Context, accountId string) (*aggregate.Account, error) {
	return a.repository.Find(ctx, accountId)
}

func NewAccountCommand() AccountCommand {
	return &accountCommand{
		repository: repository.NewAccountEventRepository(),
	}
}
