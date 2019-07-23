package query

import (
	"context"
	"github.com/payfazz/go-apt/example/eventsource/domain/account/command/aggregate"
	"github.com/payfazz/go-apt/example/eventsource/domain/account/query/model"
	"github.com/payfazz/go-apt/example/eventsource/domain/account/query/repository"
)

type AccountQuery interface {
	All(ctx context.Context) ([]*model.Account, error)
	Get(ctx context.Context, accountId string) (*model.Account, error)
	DirectUpdate(ctx context.Context, account *aggregate.Account) error
}

type accountQuery struct {
	repository repository.AccountRepository
}

func (a *accountQuery) All(ctx context.Context) ([]*model.Account, error) {
	return a.repository.All(ctx)
}

func (a *accountQuery) Get(ctx context.Context, accountId string) (*model.Account, error) {
	return a.repository.Find(ctx, accountId)
}

func (a *accountQuery) DirectUpdate(ctx context.Context, account *aggregate.Account) error {
	accModel, err := a.repository.Find(ctx, account.Id)
	if err != nil {
		return err
	}

	if account.DeletedAt != nil {
		_ = a.repository.Delete(ctx, accModel)
		return nil
	}

	newAcc := accModel
	if newAcc == nil {
		newAcc = model.AccountModel()
	}

	if newAcc.Version >= account.Version {
		// return if already up to date
		return nil
	}

	newAcc.Id = account.Id
	newAcc.Version = account.Version
	newAcc.Name = account.Name
	newAcc.Balance = account.Balance

	if accModel == nil {
		_, err = a.repository.Create(ctx, newAcc)
	} else {
		err = a.repository.Update(ctx, newAcc)
	}

	return err

}

func NewAccountQuery() AccountQuery {
	return &accountQuery{
		repository: repository.NewAccountRepository(),
	}
}
