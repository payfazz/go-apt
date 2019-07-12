package es

import (
	"github.com/jmoiron/sqlx/types"
	"time"
)

// Event is a struct for event
type Event struct {
	Id        int64          `json:"id" db:"id"`
	Type      string         `json:"type" db:"type"`
	Data      types.JSONText `json:"data" db:"data"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
}
