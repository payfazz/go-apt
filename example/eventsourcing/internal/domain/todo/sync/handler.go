package sync

import (
	"context"
	"github.com/payfazz/go-apt/example/eventsourcing/internal/domain/todo/data"
	"github.com/payfazz/go-apt/example/eventsourcing/internal/domain/todo/query"
)

// TodoSyncHandler is event handler interface for todo event
type TodoSyncHandler interface {
	HandleTodoCreated(ctx context.Context, data data.TodoCreated) error
	HandleTodoUpdated(ctx context.Context, data data.TodoUpdated) error
	HandleTodoDeleted(ctx context.Context, data data.TodoDeleted) error
}

type todoSyncHandler struct {
	repository query.TodoReadRepository
}

func (t *todoSyncHandler) HandleTodoCreated(ctx context.Context, data data.TodoCreated) error {
	todo := query.TodoReadModel()
	todo.Id = data.Id
	todo.Text = data.Text
	todo.Completed = false
	_, err := t.repository.Create(ctx, todo)
	return err
}

func (t *todoSyncHandler) HandleTodoUpdated(ctx context.Context, data data.TodoUpdated) error {
	todo := query.TodoReadModel()
	todo.Id = data.Id
	todo.Text = data.Text
	todo.Completed = data.Completed
	err := t.repository.Update(ctx, todo)
	return err
}

func (t *todoSyncHandler) HandleTodoDeleted(ctx context.Context, data data.TodoDeleted) error {
	todo := query.TodoReadModel()
	todo.Id = data.Id
	err := t.repository.Delete(ctx, todo)
	return err
}

func NewTodoSyncHandler(repository query.TodoReadRepository) TodoSyncHandler {
	return &todoSyncHandler{repository: repository}
}
