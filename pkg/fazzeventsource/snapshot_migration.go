package fazzeventsource

import "github.com/payfazz/go-apt/pkg/fazzdb"

func CreateSnapshotsTable(name string) *fazzdb.MigrationTable {
	return fazzdb.CreateTable(name, func(table *fazzdb.MigrationTable) {
		table.Field(fazzdb.CreateUuid("id").Primary())
		table.Field(fazzdb.CreateInteger("version"))
		table.Field(fazzdb.CreateJsonB("data"))
	})
}
