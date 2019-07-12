package todo

import (
	"context"
	"github.com/payfazz/go-apt/example/eventsourcing/internal/domain/todo/command"
	"github.com/payfazz/go-apt/example/eventsourcing/internal/domain/todo/data"
	"github.com/payfazz/go-apt/example/eventsourcing/internal/domain/todo/query"
)

// ServiceInterface interface that used for serving the service
type ServiceInterface interface {
	All(ctx context.Context) ([]query.Todo, error)
	Create(ctx context.Context, payload data.PayloadCreateTodo) (*string, error)
	Update(ctx context.Context, payload data.PayloadUpdateTodo) error
	Delete(ctx context.Context, id string) error
}

type service struct {
	Command command.TodoCommand
}

// All is a function to get all todo task
func (s *service) All(ctx context.Context) ([]query.Todo, error) {
	panic("implement me")
}

// Create is a function to create a todo task
func (s *service) Create(ctx context.Context, payload data.PayloadCreateTodo) (*string, error) {
	return s.Command.Create(ctx, payload)
}

// Update is a function to update a todo task
func (s *service) Update(ctx context.Context, payload data.PayloadUpdateTodo) error {
	return s.Command.Update(ctx, payload)
}

// Delete is a function to update a todo task
func (s *service) Delete(ctx context.Context, id string) error {
	return s.Command.Delete(ctx, id)
}

// NewTodoService is a function to construct todo service
func NewTodoService() ServiceInterface {
	return &service{
		Command: command.NewTodoCommand(),
	}
}
