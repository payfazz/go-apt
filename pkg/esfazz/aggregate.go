package esfazz

import (
	"encoding/json"
)

// AggregateRow is a model for aggregate snapshot in database
type AggregateRow struct {
	Id      string          `json:"id" db:"id"`
	Version int             `json:"version" db:"version"`
	Data    json.RawMessage `json:"data" db:"data"`
}

// Aggregate is interface for aggregate object
type Aggregate interface {
	GetId() string
	GetVersion() int
	Apply(eventLog *EventLog) error
}

// BaseAggregate is a struct to be used composed with aggregate object
type BaseAggregate struct {
	Id      string `json:"id"`
	Version int    `json:"version"`
}

// GetId return Id of aggregate object
func (a *BaseAggregate) GetId() string {
	return a.Id
}

// GetVersion return aggregate version of aggregate object
func (a *BaseAggregate) GetVersion() int {
	return a.Version
}

// Apply apply event to aggregate
func (a *BaseAggregate) Apply(eventLog *EventLog) error {
	return nil
}
