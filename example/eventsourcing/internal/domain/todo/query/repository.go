package query

import (
	"context"
	"github.com/payfazz/go-apt/example/eventsourcing/lib/fazzcontext"
	"github.com/payfazz/go-apt/pkg/fazzcommon/formatter"
)

// TodoReadRepository is repository of read model for todo
type TodoReadRepository interface {
	All(ctx context.Context) ([]*Todo, error)
	Create(ctx context.Context, todo *Todo) (*string, error)
	Update(ctx context.Context, todo *Todo) error
	Delete(ctx context.Context, todo *Todo) error
}

type todoReadRepository struct {
	todo *Todo
}

func (t *todoReadRepository) All(ctx context.Context) ([]*Todo, error) {
	q, err := fazzcontext.GetQuery(ctx)
	if nil != err {
		return nil, err
	}

	rows, err := q.Use(t.todo).AllCtx(ctx)
	if nil != err {
		return nil, err
	}

	results := rows.([]*Todo)

	return results, nil
}

func (t *todoReadRepository) Create(ctx context.Context, todo *Todo) (*string, error) {
	q, err := fazzcontext.GetQuery(ctx)
	if nil != err {
		return nil, err
	}

	result, err := q.Use(todo).
		InsertCtx(ctx, false)

	if nil != err {
		return nil, err
	}

	id := formatter.SliceUint8ToString(result.([]uint8))
	return &id, nil
}

func (t *todoReadRepository) Update(ctx context.Context, todo *Todo) error {
	q, err := fazzcontext.GetQuery(ctx)
	if err != nil {
		return err
	}

	_, err = q.Use(todo).
		UpdateCtx(ctx)

	return err
}

func (t *todoReadRepository) Delete(ctx context.Context, todo *Todo) error {
	q, err := fazzcontext.GetQuery(ctx)
	if nil != err {
		return err
	}

	_, err = q.Use(todo).
		DeleteCtx(ctx)

	return err
}

func NewTodoReadRepository(todo *Todo) TodoReadRepository {
	return &todoReadRepository{todo: todo}
}
