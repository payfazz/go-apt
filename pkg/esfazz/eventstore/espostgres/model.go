package espostgres

import (
	"github.com/jmoiron/sqlx/types"
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

// EventLog is struct for model in database
type EventLog struct {
	fazzdb.Model
	EventId          int64          `db:"event_id"`
	EventType        string         `db:"event_type"`
	AggregateId      string         `db:"aggregate_id"`
	AggregateVersion int            `db:"aggregate_version"`
	Data             types.JSONText `db:"data"`
}

// Get is a function that used to get the data from table
func (m *EventLog) Get(key string) interface{} {
	return m.Payload()[key]
}

// Payload is a function that used to get the payload data
func (m *EventLog) Payload() map[string]interface{} {
	return m.MapPayload(m)
}

// EventLogModel create
func EventLogModel(tableName string) *EventLog {
	return &EventLog{
		Model: fazzdb.AutoIncrementModel(tableName,
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
