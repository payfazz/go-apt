package query

import (
	"context"
)

// TodoQuery is an interface for todo query
type TodoQuery interface {
	All(ctx context.Context) ([]*Todo, error)
}

type todoQuery struct {
	repository TodoReadRepository
}

func (t *todoQuery) All(ctx context.Context) ([]*Todo, error) {
	return t.repository.All(ctx)
}

func NewTodoQuery(repository TodoReadRepository) TodoQuery {
	return &todoQuery{repository: repository}
}
