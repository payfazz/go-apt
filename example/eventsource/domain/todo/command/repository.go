package command

import (
	"context"
	"encoding/json"
	"github.com/payfazz/go-apt/pkg/fazzeventsource"
)

// TodoWriteRepository is repository for todo event
type TodoWriteRepository interface {
	Post(ctx context.Context, id string, eventType string, eventData interface{}) (*fazzeventsource.Event, error)
	Get(ctx context.Context, id string) (*Todo, error)
	IsExists(ctx context.Context, id string) (bool, error)
}

type todoWriteRepository struct {
	evStore   fazzeventsource.EventStore
	snapStore fazzeventsource.SnapshotStore
	publisher fazzeventsource.EventPublisher
}

// Post do save and post for event
func (t *todoWriteRepository) Post(ctx context.Context, id string, eventType string, eventData interface{}) (*fazzeventsource.Event, error) {
	// save to event store
	savedEvent, err := t.evStore.Save(ctx, eventType, eventData)
	if err != nil {
		return nil, err
	}

	// publish to pubsub
	err = t.publisher.Publish(ctx, savedEvent)
	if err != nil {
		return nil, err
	}

	// save snapshot
	todo, err := t.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	_, err = t.snapStore.Save(ctx, todo.Id, savedEvent.Id, todo)
	if err != nil {
		return nil, err
	}

	return savedEvent, nil
}

// Get return Todo aggregates if exists
func (t *todoWriteRepository) Get(ctx context.Context, id string) (*Todo, error) {
	snap, err := t.snapStore.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	var todo = &Todo{}
	var evs []*fazzeventsource.Event
	if snap != nil {
		err = json.Unmarshal(snap.Data, todo)
		if err != nil {
			return nil, err
		}
		evs, err = t.evStore.FindAllByKey(ctx, "id", id, snap.LastEventId)
		if err != nil {
			return nil, err
		}
	} else {
		todo.Id = id
		evs, err = t.evStore.FindAllByKey(ctx, "id", id, 0)
		if err != nil {
			return nil, err
		}
		if len(evs) == 0 {
			return nil, nil
		}
	}

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
	evStore fazzeventsource.EventStore,
	snapStore fazzeventsource.SnapshotStore,
	publisher fazzeventsource.EventPublisher,
) TodoWriteRepository {
	return &todoWriteRepository{
		evStore:   evStore,
		snapStore: snapStore,
		publisher: publisher,
	}
}
