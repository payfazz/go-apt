package command

import (
	"context"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/payfazz/go-apt/example/eventsource/domain/todo/data"
	"github.com/payfazz/go-apt/pkg/fazzeventsource"
)

// TodoCommand is a interface for todo commands
type TodoCommand interface {
	Create(ctx context.Context, payload data.PayloadCreateTodo) (*string, error)
	Update(ctx context.Context, payload data.PayloadUpdateTodo) error
	Delete(ctx context.Context, id string) error
}

type todoCommand struct {
	repository TodoWriteRepository
}

// Create is a command for Create Todo
func (t *todoCommand) Create(ctx context.Context, payload data.PayloadCreateTodo) (*string, error) {

	uuidV4, _ := uuid.NewV4()
	id := uuidV4.String()

	ev := fazzeventsource.EventPayload{
		AggregateId:      id,
		AggregateVersion: 0,
		Type:             data.EVENT_TODO_CREATED,
		Data: data.TodoCreated{
			Text: payload.Text,
		},
	}

	_, err := t.repository.Save(ctx, ev)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

// Update is a command for Update Todo
func (t *todoCommand) Update(ctx context.Context, payload data.PayloadUpdateTodo) error {
	todo, err := t.repository.Get(ctx, payload.Id)
	if err != nil {
		return err
	}

	if todo == nil || todo.DeletedAt != nil {
		return errors.New("todo not found")
	}

	ev := fazzeventsource.EventPayload{
		AggregateId:      todo.Id,
		AggregateVersion: todo.Version,
		Type:             data.EVENT_TODO_UPDATED,
		Data: data.TodoUpdated{
			Text:      payload.Text,
			Completed: payload.Completed,
		},
	}

	_, err = t.repository.Save(ctx, ev)
	return err
}

// Delete is a command for Delete Todo
func (t *todoCommand) Delete(ctx context.Context, id string) error {
	todo, err := t.repository.Get(ctx, id)
	if err != nil {
		return err
	}

	if todo == nil || todo.DeletedAt != nil {
		return errors.New("todo not found")
	}

	ev := fazzeventsource.EventPayload{
		AggregateId:      todo.Id,
		AggregateVersion: todo.Version,
		Type:             data.EVENT_TODO_DELETED,
	}

	_, err = t.repository.Save(ctx, ev)
	return err
}

// NewTodoCommand is a constructor for todo command handler
func NewTodoCommand(repository TodoWriteRepository) TodoCommand {
	return &todoCommand{
		repository: repository,
	}
}
