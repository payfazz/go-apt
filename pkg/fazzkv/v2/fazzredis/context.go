package fazzredis

import (
	"context"
	"errors"
)

type rdKeyType struct{}

var rdKey rdKeyType

// Context is a function to append redis object into context
func Context(ctx context.Context, rds Store) context.Context {
	return context.WithValue(ctx, rdKey, rds)
}

// GetFromContext is a function to get redis object from context
// Must be used after NewRedisContext
func GetFromContext(ctx context.Context) (Store, error) {
	rds := ctx.Value(rdKey)
	if nil == rds {
		return nil, errors.New("redis instance not found in context, make sure to call fazzredis.Context before calling fazzredis.GetFromContext")
	}
	return rds.(Store), nil
}
