package query

import (
	"context"
	"github.com/payfazz/go-apt/example/esfazz/account/model"
	"github.com/payfazz/go-apt/example/esfazz/account/query/repository"
)

// AccountQuery is query interface for account
type AccountQuery interface {
	All(ctx context.Context) ([]*model.Account, error)
	Get(ctx context.Context, accountId string) (*model.Account, error)
}

type accountQuery struct {
	repository repository.AccountRepository
}

// All return all account
func (a *accountQuery) All(ctx context.Context) ([]*model.Account, error) {
	return a.repository.All(ctx)
}

// Get return account by id
func (a *accountQuery) Get(ctx context.Context, accountId string) (*model.Account, error) {
	return a.repository.Find(ctx, accountId)
}

// NewAccountQuery is constructor for AccountQuery
func NewAccountQuery() AccountQuery {
	return &accountQuery{
		repository: repository.NewAccountRepository(),
	}
}
