package command

import (
	"context"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/payfazz/go-apt/example/eventsourcing/internal/domain/todo/data"
)

// TodoCommand is a interface for todo commands
type TodoCommand interface {
	Create(ctx context.Context, payload data.PayloadCreateTodo) (*string, error)
	Update(ctx context.Context, payload data.PayloadUpdateTodo) error
	Delete(ctx context.Context, id string) error
}

type todoCommand struct {
	repository TodoEventRepository
}

// Create is a command for Create Todo
func (t *todoCommand) Create(ctx context.Context, payload data.PayloadCreateTodo) (*string, error) {

	uuidV4, _ := uuid.NewV4()
	id := uuidV4.String()

	eventData := data.TodoCreated{
		Id:   id,
		Text: payload.Text,
	}
	savedEvent, err := t.repository.Save(ctx, data.EVENT_TODO_CREATED, eventData)
	if err != nil {
		return nil, err
	}

	err = t.repository.Publish(ctx, "", savedEvent)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

// Update is a command for Update Todo
func (t *todoCommand) Update(ctx context.Context, payload data.PayloadUpdateTodo) error {
	todoExists, err := t.repository.IsExists(ctx, payload.Id)
	if err != nil {
		return err
	}
	if !todoExists {
		return errors.New("todo not found")
	}

	eventData := data.TodoUpdated(payload)
	savedEvent, err := t.repository.Save(ctx, data.EVENT_TODO_UPDATED, eventData)
	if err != nil {
		return err
	}

	err = t.repository.Publish(ctx, "", savedEvent)
	return err
}

// Delete is a command for Delete Todo
func (t *todoCommand) Delete(ctx context.Context, id string) error {
	todoExists, err := t.repository.IsExists(ctx, id)
	if err != nil {
		return err
	}
	if !todoExists {
		return errors.New("todo not found")
	}

	eventData := data.TodoDeleted{Id: id}
	savedEvent, err := t.repository.Save(ctx, data.EVENT_TODO_DELETED, eventData)
	if err != nil {
		return err
	}

	err = t.repository.Publish(ctx, "", savedEvent)
	return err
}

// NewTodoCommand is a constructor for todo command handler
func NewTodoCommand(repository TodoEventRepository) TodoCommand {
	return &todoCommand{
		repository: repository,
	}
}
