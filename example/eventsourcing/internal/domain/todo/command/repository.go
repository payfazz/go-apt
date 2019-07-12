package command

import (
	"context"
	"encoding/json"
	"github.com/payfazz/go-apt/example/eventsourcing/lib/es"
	"github.com/payfazz/go-apt/example/eventsourcing/lib/pubsub"
)

// TodoEventRepository is repository for todo event
type TodoEventRepository interface {
	es.EventRepository
	IsExists(ctx context.Context, id string) (bool, error)
	Publish(ctx context.Context, event *es.Event) error
}

type todoEventRepository struct {
	es.EventRepository
	bus pubsub.PubSub
}

// IsExists check if todo exists and not deleted
func (t *todoEventRepository) IsExists(ctx context.Context, id string) (bool, error) {
	evs, err := t.AggregateById(ctx, id)
	if err != nil {
		return false, err
	}

	exists := len(evs) > 0 && evs[len(evs)-1].Type != "todo_deleted"
	return exists, nil
}

// Publish do publishing to event bus
func (t *todoEventRepository) Publish(ctx context.Context, ev *es.Event) error {
	evJson, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	err = t.bus.Publish(ctx, ev.Type, evJson)
	return err

}

// NewTodoEventRepository is constructor for Todo Event Repository
func NewTodoEventRepository() TodoEventRepository {
	return &todoEventRepository{
		EventRepository: es.NewEventStore(),
		bus:             pubsub.SingletonPubSub(),
	}
}
