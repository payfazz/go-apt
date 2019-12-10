package fazzeventsource

import "github.com/payfazz/go-apt/pkg/fazzdb"

// CreateEventsTable return migration table for event
func CreateEventsTable(tableName string) *fazzdb.MigrationTable {
	return fazzdb.CreateTable(tableName, func(table *fazzdb.MigrationTable) {
		table.Field(fazzdb.CreateBigSerial("event_id").Primary())
		table.Field(fazzdb.CreateUuid("entity_id"))
		table.Field(fazzdb.CreateInteger("entity_version"))
		table.Field(fazzdb.CreateString("event_type"))
		table.Field(fazzdb.CreateJsonB("event_data"))
		table.TimestampsTz(0)
		table.Uniques("entity_id", "entity_version")
	})
}
