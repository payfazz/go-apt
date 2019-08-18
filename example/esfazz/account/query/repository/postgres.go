package repository

import (
	"context"
	"github.com/payfazz/go-apt/example/esfazz/account/model"
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

// AccountRepository is repository interface for account
type AccountRepository interface {
	All(ctx context.Context) ([]*model.Account, error)
	Find(ctx context.Context, id string) (*model.Account, error)
}

type accountRepository struct {
	account *model.Account
}

// All return all account
func (a *accountRepository) All(ctx context.Context) ([]*model.Account, error) {
	q, err := fazzdb.GetQueryContext(ctx)
	if nil != err {
		return nil, err
	}

	rows, err := q.Use(a.account).
		WhereNil("deleted_time").
		AllCtx(ctx)
	if nil != err {
		return nil, err
	}

	results := rows.([]*model.Account)

	return results, nil
}

// Find return account by id
func (a *accountRepository) Find(ctx context.Context, id string) (*model.Account, error) {
	q, err := fazzdb.GetQueryContext(ctx)
	if nil != err {
		return nil, err
	}

	rows, err := q.Use(a.account).
		Where("id", id).
		WhereNil("deleted_time").
		WithLimit(1).
		AllCtx(ctx)

	if nil != err {
		return nil, err
	}

	results := rows.([]*model.Account)
	if len(results) == 0 {
		return nil, nil
	}

	return results[0], nil
}

// NewAccountRepository is constructor for AccountRepository
func NewAccountRepository() AccountRepository {
	return &accountRepository{account: model.AccountModel()}
}
