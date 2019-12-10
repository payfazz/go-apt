package fazzeventsource

import "github.com/payfazz/go-apt/pkg/fazzdb"

// EventLog is struct for event log structure in database
type EventLog struct {
	fazzdb.Model
	EventId       int64  `db:"event_id"`
	EntityId      string `db:"entity_id"`
	EntityVersion int    `db:"entity_version"`
	Event
}

// Get is a function that used to get the data from table
func (m *EventLog) Get(key string) interface{} {
	return m.Payload()[key]
}

// Payload is a function that used to get the payload data
func (m *EventLog) Payload() map[string]interface{} {
	return m.MapPayload(m)
}

// EventLogModel create EventLog model
func EventLogModel(tableName string) *EventLog {
	return &EventLog{
		Model: fazzdb.AutoIncrementModel(tableName,
			[]fazzdb.Column{
				fazzdb.Col("event_id"),
				fazzdb.Col("entity_id"),
				fazzdb.Col("entity_version"),
				fazzdb.Col("event_type"),
				fazzdb.Col("event_data"),
			},
			"event_id",
			true,
			false,
		),
	}
}
