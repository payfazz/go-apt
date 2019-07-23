package repository

import (
	"context"
	"github.com/payfazz/go-apt/example/eventsource/domain/account/query/model"
	"github.com/payfazz/go-apt/pkg/fazzcommon/formatter"
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

type AccountRepository interface {
	All(ctx context.Context) ([]*model.Account, error)
	Find(ctx context.Context, id string) (*model.Account, error)
	Create(ctx context.Context, todo *model.Account) (*string, error)
	Update(ctx context.Context, todo *model.Account) error
	Delete(ctx context.Context, todo *model.Account) error
}

type accountRepository struct {
	account *model.Account
}

func (a *accountRepository) All(ctx context.Context) ([]*model.Account, error) {
	q, err := fazzdb.GetQueryContext(ctx)
	if nil != err {
		return nil, err
	}

	rows, err := q.Use(a.account).AllCtx(ctx)
	if nil != err {
		return nil, err
	}

	results := rows.([]*model.Account)

	return results, nil
}

func (a *accountRepository) Find(ctx context.Context, id string) (*model.Account, error) {
	q, err := fazzdb.GetQueryContext(ctx)
	if nil != err {
		return nil, err
	}

	rows, err := q.Use(a.account).
		Where("id", id).
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

func (a *accountRepository) Create(ctx context.Context, account *model.Account) (*string, error) {
	q, err := fazzdb.GetQueryContext(ctx)
	if nil != err {
		return nil, err
	}

	result, err := q.Use(account).InsertCtx(ctx, false)

	if nil != err {
		return nil, err
	}

	id := formatter.SliceUint8ToString(result.([]uint8))
	return &id, nil
}

func (a *accountRepository) Update(ctx context.Context, account *model.Account) error {
	q, err := fazzdb.GetQueryContext(ctx)
	if err != nil {
		return err
	}

	_, err = q.Use(account).UpdateCtx(ctx)

	return err
}

func (a *accountRepository) Delete(ctx context.Context, account *model.Account) error {
	q, err := fazzdb.GetQueryContext(ctx)
	if nil != err {
		return err
	}

	_, err = q.Use(account).DeleteCtx(ctx)

	return err
}

func NewAccountRepository() AccountRepository {
	return &accountRepository{account: model.AccountModel()}
}
