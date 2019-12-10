package fazzeventsource

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Event is struct for event domain
type Event struct {
	Type string    `db:"event_type" json:"event_type"`
	Data EventData `db:"event_data" json:"event_data"`
}

// EventData type wrapper for event data map
type EventData map[string]interface{}

// Value implement value for sql
func (e EventData) Value() (driver.Value, error) {
	return json.Marshal(e)
}

// Scan implement scanner for sql data
func (e *EventData) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &e)
}
