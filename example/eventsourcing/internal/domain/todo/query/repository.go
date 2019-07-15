package query

import "context"

// TodoReadRepository is repository of read model for todo
type TodoReadRepository interface {
	All(ctx context.Context) error
}
