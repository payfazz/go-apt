package command

import (
	"context"
	"github.com/payfazz/go-apt/pkg/fazzeventsource"
)

// TodoWriteRepository is repository for todo event
type TodoWriteRepository interface {
	fazzeventsource.EventStore
	fazzeventsource.EventPublisher
	Get(ctx context.Context, id string) (*Todo, error)
	IsExists(ctx context.Context, id string) (bool, error)
}

type todoWriteRepository struct {
	fazzeventsource.EventStore
	fazzeventsource.EventPublisher
}

// Get return Todo aggregates if exists
func (t *todoWriteRepository) Get(ctx context.Context, id string) (*Todo, error) {
	evs, err := t.FindByInstanceId(ctx, id)
	if err != nil {
		return nil, err
	}

	todo := &Todo{Id: id}
	err = todo.ApplyAll(evs...)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

// IsExists check if todo exists and not deleted
func (t *todoWriteRepository) IsExists(ctx context.Context, id string) (bool, error) {
	todo, err := t.Get(ctx, id)
	if err != nil {
		return false, err
	}
	return todo.DeletedAt == nil, nil
}

// NewTodoEventRepository is constructor for Todo Event Repository
func NewTodoEventRepository(store fazzeventsource.EventStore, publisher fazzeventsource.EventPublisher) TodoWriteRepository {
	return &todoWriteRepository{
		EventStore:     store,
		EventPublisher: publisher,
	}
}
