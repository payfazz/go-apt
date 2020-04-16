package redis

import (
	"errors"
	"time"

	"github.com/go-redis/redis"
	"github.com/payfazz/go-apt/pkg/fazzkv"
)

// RedisInterface is abstraction layer redis that wrap store interface with addition
// for adding expire time in redis set.
type RedisInterface interface {
	fazzkv.Store
	GetClient() *redis.Client
	Increment(key string) error
	SetWithExpire(key string, value interface{}, duration time.Duration) error
	SetWithExpireIfNotExist(key string, value interface{}, duration time.Duration) error
}

// private struct for wrapping go-redis client
type fazzRedis struct {
	client *redis.Client
}

// Set accept key (string) and value (any), return error if it's failed to set the data,
// this method allow user to save the data to redis with K-V mechanism.
func (kv *fazzRedis) Set(key string, value interface{}) error {
	return kv.client.Set(key, value, 0).Err()
}

// Get accept key (string)  and return error if it's failed to get the data,
// this method allow user to get the data from redis with K-V mechanism.
func (kv *fazzRedis) Get(key string) (string, error) {
	return kv.client.Get(key).Result()
}

// Delete accept key (string) return error if it's failed to delete the data.
func (kv *fazzRedis) Delete(key string) error {
	return kv.client.Del(key).Err()
}

// Truncate allow user to remove all data from redis.
func (kv *fazzRedis) Truncate() error {
	return kv.client.FlushAll().Err()
}

// Increment allow user to increment integer data without resetting expiry time
func (kv *fazzRedis) Increment(key string) error {
	return kv.client.Incr(key).Err()
}

// SetWithExpire allow user to set data and expired time at one time.
func (kv *fazzRedis) SetWithExpire(key string, value interface{}, duration time.Duration) error {
	return kv.client.Set(key, value, duration).Err()
}

// SetWithExpireIfNotExists allow user to set data and expired time at one time.
// It returns error if key already exists
func (kv *fazzRedis) SetWithExpireIfNotExist(key string, value interface{}, duration time.Duration) error {
	set, err := kv.client.SetNX(key, value, duration).Result()
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

// NewFazzRedis is a function that act as constructor and injector for FazzRedis.
func NewFazzRedis(addr string, password string) (RedisInterface, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	if err := client.Ping().Err(); err != nil {
		return nil, err
	}
	return &fazzRedis{client: client}, nil
}
