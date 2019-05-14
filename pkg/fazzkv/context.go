package fazzkv

import (
	"context"

	"github.com/payfazz/go-apt/pkg/fazzkv/redis"
)

type rdKeyType struct{}

var rdKey rdKeyType

// NewRedisContext is a function to append redis object into context
func NewRedisContext(ctx context.Context, addr string, password string) context.Context {
	rds, _ := redis.NewFazzRedis(addr, password)
	return context.WithValue(ctx, rdKey, rds)
}

// GetRedisContext is a function to get redis object from context
// Must be used after NewRedisContext
func GetRedisContext(ctx context.Context) redis.RedisInterface {
	return ctx.Value(rdKey).(redis.RedisInterface)
}
