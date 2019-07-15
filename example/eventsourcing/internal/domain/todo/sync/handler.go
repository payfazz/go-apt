package sync

import (
	"context"
	"github.com/payfazz/go-apt/example/eventsourcing/internal/domain/todo/data"
)

// TodoSyncEventHandler is event handler interface for todo event
type TodoSyncEventHandler interface {
	HandleTodoCreated(ctx context.Context, data data.TodoCreated) error
	HandleTodoUpdated(ctx context.Context, data data.TodoUpdated) error
	HandleTodoDeleted(ctx context.Context, data data.TodoDeleted) error
}
