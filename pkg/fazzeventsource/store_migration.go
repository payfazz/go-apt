package fazzeventsource

import "github.com/payfazz/go-apt/pkg/fazzdb"

func CreateEventsTable(name string) *fazzdb.MigrationTable {
	return fazzdb.CreateTable(name, func(table *fazzdb.MigrationTable) {
		table.Field(fazzdb.CreateBigSerial("id").Primary())
		table.Field(fazzdb.CreateString("type"))
		table.Field(fazzdb.CreateJsonB("data"))
		table.Field(fazzdb.CreateTimestampTz("created_at", 0))
	})
}
