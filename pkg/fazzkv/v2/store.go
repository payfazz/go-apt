package fazzkv

import "context"

// Store is an abstraction of key value storing data
// it can be used for storing data to redis, memcached, etc.
type Store interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Truncate(ctx context.Context) error
}
