package query

import (
	"context"
	"github.com/payfazz/go-apt/example/fazzeventsource_sample/internal/domain/todo/data"
)

// TodoSyncHandler is event handler interface for todo event
type TodoSyncHandler interface {
	HandleTodoCreated(ctx context.Context, data data.TodoCreated) error
	HandleTodoUpdated(ctx context.Context, data data.TodoUpdated) error
	HandleTodoDeleted(ctx context.Context, data data.TodoDeleted) error
}

type todoSyncHandler struct {
	repository TodoReadRepository
}

func (t *todoSyncHandler) HandleTodoCreated(ctx context.Context, data data.TodoCreated) error {
	todo := TodoReadModel()
	todo.Id = data.Id
	todo.Text = data.Text
	todo.Completed = false
	_, err := t.repository.Create(ctx, todo)
	return err
}

func (t *todoSyncHandler) HandleTodoUpdated(ctx context.Context, data data.TodoUpdated) error {
	todo := TodoReadModel()
	todo.Id = data.Id
	todo.Text = data.Text
	todo.Completed = data.Completed
	err := t.repository.Update(ctx, todo)
	return err
}

func (t *todoSyncHandler) HandleTodoDeleted(ctx context.Context, data data.TodoDeleted) error {
	todo := TodoReadModel()
	todo.Id = data.Id
	err := t.repository.Delete(ctx, todo)
	return err
}

func NewTodoSyncHandler(repository TodoReadRepository) TodoSyncHandler {
	return &todoSyncHandler{repository: repository}
}
