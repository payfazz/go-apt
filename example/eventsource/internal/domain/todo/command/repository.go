package command

import (
	"context"
	"github.com/payfazz/go-apt/pkg/fazzeventsource"
)

// TodoWriteRepository is repository for todo event
type TodoWriteRepository interface {
	Post(ctx context.Context, eventType string, eventData interface{}) (*fazzeventsource.Event, error)
	Get(ctx context.Context, id string) (*Todo, error)
	IsExists(ctx context.Context, id string) (bool, error)
}

type todoWriteRepository struct {
	store     fazzeventsource.EventStore
	publisher fazzeventsource.EventPublisher
}

// Post do save and post for event
func (t *todoWriteRepository) Post(ctx context.Context, eventType string, eventData interface{}) (*fazzeventsource.Event, error) {
	savedEvent, err := t.store.Save(ctx, eventType, eventData)
	if err != nil {
		return nil, err
	}

	err = t.publisher.Publish(ctx, savedEvent)

	if err != nil {
		return nil, err
	}
	return savedEvent, nil
}

// Get return Todo aggregates if exists
func (t *todoWriteRepository) Get(ctx context.Context, id string) (*Todo, error) {
	evs, err := t.store.FindByInstanceId(ctx, id)
	if err != nil {
		return nil, err
	}
	if len(evs) == 0 {
		return nil, nil
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
func NewTodoEventRepository(
	store fazzeventsource.EventStore,
	publisher fazzeventsource.EventPublisher,
) TodoWriteRepository {
	return &todoWriteRepository{
		store:     store,
		publisher: publisher,
	}
}
