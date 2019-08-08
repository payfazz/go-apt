package eventpostgres

import (
	"github.com/jmoiron/sqlx/types"
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

// eventLog is struct for model in database
type eventLog struct {
	fazzdb.Model
	EventId          int64          `db:"event_id"`
	EventType        string         `db:"event_type"`
	AggregateId      string         `db:"aggregate_id"`
	AggregateVersion int            `db:"aggregate_version"`
	Data             types.JSONText `db:"data"`
}

// Get is a function that used to get the data from table
func (m *eventLog) Get(key string) interface{} {
	return m.Payload()[key]
}

// Payload is a function that used to get the payload data
func (m *eventLog) Payload() map[string]interface{} {
	return m.MapPayload(m)
}

// EventLogModel create
func EventLogModel(tableName string) *eventLog {
	return &eventLog{
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
