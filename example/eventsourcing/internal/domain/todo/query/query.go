package query

import (
	"context"
)

// TodoQuery is an interface for todo query
type TodoQuery interface {
	All(ctx context.Context) ([]Todo, error)
}
