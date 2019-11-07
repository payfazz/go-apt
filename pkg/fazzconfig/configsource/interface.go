package configsource

// ConfigSource is interface for config source
type ConfigSource interface {
	Get(key string) string
	Set(key string, value string) error
}
