package esfazz

import (
	"encoding/json"
)

// Event is struct for event
type Event struct {
	Type      string
	Aggregate Aggregate
	Data      json.RawMessage
}

// EventPayload is payload for creating event
type EventPayload struct {
	Type      string
	Aggregate Aggregate
	Data      interface{}
}

// CreateEvent create event from payload
func CreateEvent(payload *EventPayload) (*Event, error) {
	data, err := json.Marshal(payload.Data)

	if err != nil {
		return nil, err
	}

	return &Event{
		Type:      payload.Type,
		Aggregate: payload.Aggregate,
		Data:      data,
	}, nil
}
