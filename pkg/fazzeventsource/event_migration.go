package fazzeventsource

import "github.com/payfazz/go-apt/pkg/fazzdb"

func CreateEventsTable(name string) *fazzdb.MigrationTable {
	return fazzdb.CreateTable(name, func(table *fazzdb.MigrationTable) {
		table.Field(fazzdb.CreateBigSerial("event_id").Primary())
		table.Field(fazzdb.CreateString("event_type"))
		table.Field(fazzdb.CreateUuid("aggregate_id"))
		table.Field(fazzdb.CreateInteger("aggregate_version"))
		table.Field(fazzdb.CreateJsonB("data"))
		table.Field(fazzdb.CreateTimestampTz("created_at", 0))
	})
}
