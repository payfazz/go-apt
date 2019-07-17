package fazzeventsource

import (
	"encoding/json"
	"time"
)

// Event is a struct for event
type Event struct {
	Id        int64           `json:"id" db:"id"`
	Type      string          `json:"type" db:"type"`
	Data      json.RawMessage `json:"data" db:"data"`
	CreatedAt *time.Time      `json:"created_at" db:"created_at"`
}
