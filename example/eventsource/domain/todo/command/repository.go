package command

import (
	"context"
	"encoding/json"
	"github.com/payfazz/go-apt/pkg/fazzeventsource"
)

// TodoWriteRepository is repository for todo event
type TodoWriteRepository interface {
	Save(ctx context.Context, payload fazzeventsource.EventPayload) (*fazzeventsource.EventLog, error)
	Get(ctx context.Context, id string) (*Todo, error)
}

type todoWriteRepository struct {
	evStore   fazzeventsource.EventStore
	snapStore fazzeventsource.SnapshotStore
}

// Save do save and post for event
func (t *todoWriteRepository) Save(ctx context.Context, payload fazzeventsource.EventPayload) (*fazzeventsource.EventLog, error) {
	// save to event store
	savedEvent, err := t.evStore.Save(ctx, payload)
	if err != nil {
		return nil, err
	}

	// save snapshot
	todo, err := t.Get(ctx, savedEvent.AggregateId)
	if err != nil {
		return nil, err
	}
	_, err = t.snapStore.Save(ctx, savedEvent, todo)
	if err != nil {
		return nil, err
	}

	return savedEvent, nil
}

// Get return latest todo aggregates if exists
func (t *todoWriteRepository) Get(ctx context.Context, id string) (*Todo, error) {
	snap, err := t.snapStore.FindBy(ctx, id)
	if err != nil {
		return nil, err
	}

	todo := &Todo{}

	// load data
	if snap != nil {
		err = json.Unmarshal(snap.Data, todo)
		if err != nil {
			return nil, err
		}
	}

	evs, err := t.evStore.FindAllBy(ctx, id, todo.Version)
	if err != nil {
		return nil, err
	}

	err = todo.ApplyAll(evs...)
	if err != nil {
		return nil, err
	}

	if todo.Version == 0 {
		return nil, nil
	}

	return todo, nil
}

// NewTodoEventRepository is constructor for Todo EventLog Repository
func NewTodoEventRepository(
	evStore fazzeventsource.EventStore,
	snapStore fazzeventsource.SnapshotStore,
) TodoWriteRepository {
	return &todoWriteRepository{
		evStore:   evStore,
		snapStore: snapStore,
	}
}
