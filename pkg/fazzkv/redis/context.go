package redis

import (
	"context"
	"errors"
)

type rdKeyType struct{}

var rdKey rdKeyType

// NewRedisContext is a function to append redis object into context
func NewRedisContext(ctx context.Context, rds RedisInterface) context.Context {
	return context.WithValue(ctx, rdKey, rds)
}

// GetRedisContext is a function to get redis object from context
// Must be used after NewRedisContext
func GetRedisContext(ctx context.Context) (RedisInterface, error) {
	rds := ctx.Value(rdKey)
	if nil == rds {
		return nil, errors.New("redis instance not found in context, make sure to call NewRedisContext before calling GetRedisContext")
	}
	return rds.(RedisInterface), nil
}
