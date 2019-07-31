package espostgres

import "github.com/payfazz/go-apt/pkg/fazzdb"

// CreateEventsTable return migration table for event
func CreateEventsTable(name string) *fazzdb.MigrationTable {
	return fazzdb.CreateTable(name, func(table *fazzdb.MigrationTable) {
		table.Field(fazzdb.CreateBigSerial("event_id").Primary())
		table.Field(fazzdb.CreateString("event_type"))
		table.Field(fazzdb.CreateUuid("aggregate_id"))
		table.Field(fazzdb.CreateInteger("aggregate_version"))
		table.Field(fazzdb.CreateJsonB("data"))
		table.TimestampsTz(0)
		table.Uniques("aggregate_id", "aggregate_version")
	})
}
