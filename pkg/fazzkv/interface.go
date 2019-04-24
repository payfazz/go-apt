package fazzkv

// Store is an abstraction of key value storing data
// it can be used for storing data to redis, memcached, etc.
type Store interface {
	Set(key string, value interface{}) error
	Get(key string) (string, error)
	Delete(key string) error
	Truncate() error
	Increment(key string) error
}
