package configsource

import "os"

type envSource struct {
}

// Get return config from environment variable
func (e *envSource) Get(key string) string {
	return os.Getenv(key)
}

// Set set environment variable
func (e *envSource) Set(key string, value string) error {
	return os.Setenv(key, value)
}

// FromEnv return config source for environment variable
func FromEnv() ConfigSource {
	return &envSource{}
}
