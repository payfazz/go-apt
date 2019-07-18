package migration

import (
	"github.com/payfazz/go-apt/pkg/fazzdb"
	"github.com/payfazz/go-apt/pkg/fazzeventsource"
)

var Version1 = fazzdb.MigrationVersion{
	Tables: []*fazzdb.MigrationTable{
		fazzeventsource.CreateEventsTable("todo_events"),
		fazzeventsource.CreateSnapshotsTable("todo_snapshots"),
		fazzdb.CreateTable("todo_read", func(table *fazzdb.MigrationTable) {
			table.Field(fazzdb.CreateUuid("id").Primary())
			table.Field(fazzdb.CreateString("text"))
			table.Field(fazzdb.CreateBoolean("completed"))
			table.SoftDeleteTz(0)
		}),
	},
}
