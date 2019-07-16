package todo

import (
	"context"
	"github.com/payfazz/go-apt/example/fazzeventsource_sample/internal/domain/todo/command"
	"github.com/payfazz/go-apt/example/fazzeventsource_sample/internal/domain/todo/data"
	"github.com/payfazz/go-apt/example/fazzeventsource_sample/internal/domain/todo/query"
)

// ServiceInterface interface that used for serving the service
type ServiceInterface interface {
	All(ctx context.Context) ([]*query.Todo, error)
	Create(ctx context.Context, payload data.PayloadCreateTodo) (*string, error)
	Update(ctx context.Context, payload data.PayloadUpdateTodo) error
	Delete(ctx context.Context, id string) error
}

type service struct {
	command command.TodoCommand
	query   query.TodoQuery
}

// All is a function to get all todo task
func (s *service) All(ctx context.Context) ([]*query.Todo, error) {
	return s.query.All(ctx)
}

// Create is a function to create a todo task
func (s *service) Create(ctx context.Context, payload data.PayloadCreateTodo) (*string, error) {
	return s.command.Create(ctx, payload)
}

// Update is a function to update a todo task
func (s *service) Update(ctx context.Context, payload data.PayloadUpdateTodo) error {
	return s.command.Update(ctx, payload)
}

// Delete is a function to update a todo task
func (s *service) Delete(ctx context.Context, id string) error {
	return s.command.Delete(ctx, id)
}

// NewTodoService is a function to construct todo service
func NewTodoService(command command.TodoCommand, query query.TodoQuery) ServiceInterface {
	return &service{
		command: command,
		query:   query,
	}
}
