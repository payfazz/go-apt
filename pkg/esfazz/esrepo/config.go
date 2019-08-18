package esrepo

import (
	"github.com/payfazz/go-apt/pkg/esfazz"
	"github.com/payfazz/go-apt/pkg/esfazz/eventstore"
	"github.com/payfazz/go-apt/pkg/esfazz/snapstore"
)

// Config is config for repository constructor
type Config struct {
	eventStore       eventstore.EventStore
	snapshotStore    snapstore.SnapshotStore
	aggregateFactory esfazz.AggregateFactory
	listeners        []EventListener
}

// SetAggregateFactory set aggregate factory for config
func (c Config) SetAggregateFactory(factory esfazz.AggregateFactory) Config {
	c.aggregateFactory = factory
	return c
}

// SetEventStore set event store in config with mongoDB
func (c Config) SetEventStore(store eventstore.EventStore) Config {
	c.eventStore = store
	return c
}

// SetSnapshotStore set snapshot store in config with mongoDB
func (c Config) SetSnapshotStore(store snapstore.SnapshotStore) Config {
	c.snapshotStore = store
	return c
}

// AddEventListener add event listener to repository config
func (c Config) AddEventListener(listener ...EventListener) Config {
	c.listeners = append(c.listeners, listener...)
	return c
}
