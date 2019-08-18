package esfazz

import "github.com/payfazz/go-apt/pkg/fazzdb"

// Aggregate is interface for aggregate that can apply event
type Aggregate interface {
	GetId() string
	GetVersion() int64
	Apply(events ...*Event) error
}

// AggregateModel is an aggregate that is also a fazzdb model
type AggregateModel interface {
	fazzdb.ModelInterface
	Aggregate
}

// AggregateFactory is function type that create aggregate
type AggregateFactory func(id string) Aggregate

// BaseAggregate is base aggregate
type BaseAggregate struct {
	Id      string `json:"id" db:"id"`
	Version int64  `json:"version" db:"version"`
}

// GetId return id of base aggregate
func (a *BaseAggregate) GetId() string {
	return a.Id
}

// GetVersion return version of base aggregate
func (a *BaseAggregate) GetVersion() int64 {
	return a.Version
}

// Apply implemented in base aggregate to for aggregate interface
func (a *BaseAggregate) Apply(events ...*Event) error {
	return nil
}
