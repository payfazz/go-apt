package esfazz

import (
	"encoding/json"
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

// EventLog is a struct for event
type EventLog struct {
	fazzdb.Model
	EventId          int64           `json:"event_id" db:"event_id"`
	EventType        string          `json:"event_type" db:"event_type"`
	AggregateId      string          `json:"aggregate_id" db:"aggregate_id"`
	AggregateVersion int             `json:"aggregate_version" db:"aggregate_version"`
	Data             json.RawMessage `json:"data" db:"data"`
}

// Get is a function that used to get the data from table
func (m *EventLog) Get(key string) interface{} {
	return m.Payload()[key]
}

// Payload is a function that used to get the payload data
func (m *EventLog) Payload() map[string]interface{} {
	return m.MapPayload(m)
}

// EventLogModel is the constructor for event log model
func EventLogModel(table string) *EventLog {
	return &EventLog{
		Model: fazzdb.AutoIncrementModel(table,
			[]fazzdb.Column{
				fazzdb.Col("event_id"),
				fazzdb.Col("event_type"),
				fazzdb.Col("aggregate_id"),
				fazzdb.Col("aggregate_version"),
				fazzdb.Col("data"),
			},
			"event_id",
			true,
			false,
		),
	}
}

// EventPayload is a payload for event store function
type EventPayload struct {
	Type      string
	Aggregate Aggregate
	Data      interface{}
}
