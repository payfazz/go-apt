package esfazz

import (
	"encoding/json"
)

// Event is struct for event
type Event struct {
	Type      string
	Aggregate *BaseAggregate
	Data      json.RawMessage
}

// EventPayload is payload for creating event
type EventPayload struct {
	Type string
	Data interface{}
}
