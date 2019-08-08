package snappostgres

import "github.com/payfazz/go-apt/pkg/fazzdb"

// CreateSnapshotsTable return migration table for event
func CreateSnapshotsTable(name string) *fazzdb.MigrationTable {
	return fazzdb.CreateTable(name, func(table *fazzdb.MigrationTable) {
		table.Field(fazzdb.CreateUuid("id").Primary())
		table.Field(fazzdb.CreateJsonB("data"))
	})
}
