package command

import (
	"context"
	"github.com/payfazz/go-apt/example/eventsourcing/lib/fazzeventsource"
)

// TodoEventRepository is repository for todo event
type TodoEventRepository interface {
	fazzeventsource.EventStore
	IsExists(ctx context.Context, id string) (bool, error)
}

type todoEventRepository struct {
	fazzeventsource.EventStore
}

// IsExists check if todo exists and not deleted
func (t *todoEventRepository) IsExists(ctx context.Context, id string) (bool, error) {
	evs, err := t.FindByInstanceId(ctx, id)
	if err != nil {
		return false, err
	}

	exists := len(evs) > 0 && evs[len(evs)-1].Type != "todo_deleted"
	return exists, nil
}

// NewTodoEventRepository is constructor for Todo Event Repository
func NewTodoEventRepository(store fazzeventsource.EventStore) TodoEventRepository {
	return &todoEventRepository{
		EventStore: store,
	}
}
