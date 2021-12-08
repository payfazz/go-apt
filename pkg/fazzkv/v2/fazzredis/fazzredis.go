package fazzredis

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/payfazz/go-apt/pkg/fazzkv/v2"
)

// Store is abstraction layer redis that wrap store interface with addition
// for adding expire time in redis set.
type Store interface {
	fazzkv.Store
	GetClient() *redis.Client
	Keys(ctx context.Context, pattern string) ([]string, error)
	Increment(ctx context.Context, key string) (int64, error)
	SetWithExpire(ctx context.Context, key string, value interface{}, duration time.Duration) error
	SetWithExpireIfNotExist(ctx context.Context, key string, value interface{}, duration time.Duration) error
}

// private struct for wrapping go-redis client
type fazzRedis struct {
	client *redis.Client
}

// Set accept key (string) and value (any), return error if it's failed to set the data,
// this method allow user to save the data to redis with K-V mechanism.
func (kv *fazzRedis) Set(ctx context.Context, key string, value interface{}) error {
	return kv.client.Set(ctx, key, value, 0).Err()
}

// Get accept key (string)  and return error if it's failed to get the data,
// this method allow user to get the data from redis with K-V mechanism.
func (kv *fazzRedis) Get(ctx context.Context, key string) (string, error) {
	return kv.client.Get(ctx, key).Result()
}

// Delete accept key (string) return error if it's failed to delete the data.
func (kv *fazzRedis) Delete(ctx context.Context, key string) error {
	return kv.client.Del(ctx, key).Err()
}

// Truncate allow user to remove all data from redis.
func (kv *fazzRedis) Truncate(ctx context.Context) error {
	return kv.client.FlushAll(ctx).Err()
}

// Keys allow user to get all keys by pattern from redis.
func (kv *fazzRedis) Keys(ctx context.Context, pattern string) ([]string, error) {
	return kv.client.Keys(ctx, pattern).Result()
}

// Increment allow user to increment integer data without resetting expiry time
func (kv *fazzRedis) Increment(ctx context.Context, key string) (int64, error) {
	return kv.client.Incr(ctx, key).Result()
}

// SetWithExpire allow user to set data and expired time at one time.
func (kv *fazzRedis) SetWithExpire(ctx context.Context, key string, value interface{}, duration time.Duration) error {
	return kv.client.Set(ctx, key, value, duration).Err()
}

// SetWithExpireIfNotExist allow user to set data and expired time at one time.
// It returns error if key already exists
func (kv *fazzRedis) SetWithExpireIfNotExist(ctx context.Context, key string, value interface{}, duration time.Duration) error {
	set, err := kv.client.SetNX(ctx, key, value, duration).Result()
	if err != nil {
		return err
	}
	if !set {
		return errors.New("key exists")
	}
	return nil
}

// GetClient returns the underlying redis client connection
func (kv *fazzRedis) GetClient() *redis.Client {
	return kv.client
}

// New is a function that act as constructor and injector for FazzRedis.
func New(client *redis.Client) (Store, error) {
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return &fazzRedis{client: client}, nil
}
