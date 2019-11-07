package configsource

import "sync"

type mapSource struct {
	mu   sync.RWMutex
	data map[string]string
}

// Get return config from map
func (m *mapSource) Get(key string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.data[key]
}

// Set return config from map
func (m *mapSource) Set(key string, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
	return nil
}

// FromMap return config source for map variable
func FromMap(configMap map[string]string) ConfigSource {
	data := make(map[string]string)
	for k, v := range configMap {
		data[k] = v
	}
	return &mapSource{data: data}
}
