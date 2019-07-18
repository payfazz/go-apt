package fazzeventsource

import "github.com/payfazz/go-apt/pkg/fazzdb"

func CreateSnapshotsTable(name string) *fazzdb.MigrationTable {
	return fazzdb.CreateTable(name, func(table *fazzdb.MigrationTable) {
		table.Field(fazzdb.CreateUuid("aggregate_id").Primary())
		table.Field(fazzdb.CreateBigInteger("last_event_id"))
		table.Field(fazzdb.CreateJsonB("data"))
	})
}
