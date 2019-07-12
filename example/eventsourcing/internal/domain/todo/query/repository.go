package query

import "context"

type TodoReadRepository interface {
	All(ctx context.Context) error
}
