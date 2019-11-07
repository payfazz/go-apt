package fazzconfig

import "github.com/payfazz/go-apt/pkg/fazzconfig/configsource"

// ConfigReader is config reader interface
type ConfigReader interface {
	Get(key string) string
}

type configReader struct {
	sources []configsource.ConfigSource
}

// Get return config for a key
func (c *configReader) Get(key string) string {
	for i := range c.sources {
		result := c.sources[i].Get(key)
		if result != "" {
			return result
		}
	}
	return ""
}

// NewReader return new config reader
func NewReader(sources ...configsource.ConfigSource) ConfigReader {
	return &configReader{sources: sources}
}
