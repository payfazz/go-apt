package esrepo

import (
	"github.com/payfazz/go-apt/pkg/esfazz"
	"github.com/payfazz/go-apt/pkg/esfazz/eventstore"
	"github.com/payfazz/go-apt/pkg/esfazz/eventstore/eventmongo"
	"github.com/payfazz/go-apt/pkg/esfazz/eventstore/eventpostgres"
	"github.com/payfazz/go-apt/pkg/esfazz/snapstore"
	"github.com/payfazz/go-apt/pkg/esfazz/snapstore/snapmongo"
	"github.com/payfazz/go-apt/pkg/esfazz/snapstore/snappostgres"
	"go.mongodb.org/mongo-driver/mongo"
)

// RepositoryConfig is config for repository constructor
type RepositoryConfig struct {
	EventStore       eventstore.EventStore
	SnapshotStore    snapstore.SnapshotStore
	AggregateFactory esfazz.AggregateFactory
	SnapshotEventMax int
}

// SetAggregateFactory set aggregate factory for config
func (r *RepositoryConfig) SetAggregateFactory(factory esfazz.AggregateFactory) *RepositoryConfig {
	r.AggregateFactory = factory
	return r
}

// SetPostgresEventStore set event store in config with postgreSQL
func (r *RepositoryConfig) SetPostgresEventStore(tableName string) *RepositoryConfig {
	r.EventStore = eventpostgres.EventStore(tableName)
	return r
}

// SetMongoEventStore set event store in config with mongoDB
func (r *RepositoryConfig) SetMongoEventStore(collection *mongo.Collection) *RepositoryConfig {
	r.EventStore = eventmongo.EventStore(collection)
	return r
}

// SetEventStore set event store in config with mongoDB
func (r *RepositoryConfig) SetEventStore(store eventstore.EventStore) *RepositoryConfig {
	r.EventStore = store
	return r
}

// SetPostgresSnapshotStore set snapshot store in config with postgreSQL
func (r *RepositoryConfig) SetPostgresSnapshotStore(tableName string) *RepositoryConfig {
	r.SnapshotStore = snappostgres.SnapshotStore(tableName)
	return r
}

// SetMongoSnapshotStore set snapshot store in config with mongoDB
func (r *RepositoryConfig) SetMongoSnapshotStore(collection *mongo.Collection) *RepositoryConfig {
	r.SnapshotStore = snapmongo.SnapshotStore(collection)
	return r
}

// SetSnapshotStore set snapshot store in config with mongoDB
func (r *RepositoryConfig) SetSnapshotStore(store snapstore.SnapshotStore) *RepositoryConfig {
	r.SnapshotStore = store
	return r
}

// SetEventMax set snapshot event max
func (r *RepositoryConfig) SetEventMax(count int) *RepositoryConfig {
	r.SnapshotEventMax = count
	return r
}
