package event

import (
	"context"
)

// TodoEventHandler is event handler interface for todo event
type TodoEventHandler interface {
	HandleTodoCreated(ctx context.Context, data TodoCreatedData) error
	HandleTodoUpdated(ctx context.Context, data TodoUpdatedData) error
	HandleTodoDeleted(ctx context.Context, data TodoDeletedData) error
}
