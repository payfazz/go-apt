package esfazz

import (
	"encoding/json"
	"time"
)

// EventLog is a struct for event
type EventLog struct {
	EventId          int64           `json:"event_id" db:"event_id"`
	EventType        string          `json:"event_type" db:"event_type"`
	AggregateId      string          `json:"aggregate_id" db:"aggregate_id"`
	AggregateVersion int             `json:"aggregate_version" db:"aggregate_version"`
	Data             json.RawMessage `json:"data" db:"data"`
	CreatedAt        *time.Time      `json:"created_at" db:"created_at"`
}

// EventPayload is a payload for event store function
type EventPayload struct {
	Type      string
	Aggregate Aggregate
	Data      interface{}
}
