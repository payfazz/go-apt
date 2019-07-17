package command

import (
	"context"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/payfazz/go-apt/example/eventsource/internal/domain/todo/data"
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

	eventData := data.TodoCreated{
		Id:   id,
		Text: payload.Text,
	}
	_, err := t.repository.Post(ctx, data.EVENT_TODO_CREATED, eventData)
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

	eventData := data.TodoUpdated(payload)
	_, err = t.repository.Post(ctx, data.EVENT_TODO_UPDATED, eventData)
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

	eventData := data.TodoDeleted{Id: id}
	_, err = t.repository.Post(ctx, data.EVENT_TODO_DELETED, eventData)
	return err
}

// NewTodoCommand is a constructor for todo command handler
func NewTodoCommand(repository TodoWriteRepository) TodoCommand {
	return &todoCommand{
		repository: repository,
	}
}
