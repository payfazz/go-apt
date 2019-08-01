package migration

import (
	"github.com/payfazz/go-apt/pkg/esfazz/eventstore/eventpostgres"
	"github.com/payfazz/go-apt/pkg/esfazz/snapstore/snappostgres"
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

var Version1 = fazzdb.MigrationVersion{
	Tables: []*fazzdb.MigrationTable{
		eventpostgres.CreateEventsTable("account_event"),
		snappostgres.CreateSnapshotsTable("account_snap"),
		fazzdb.CreateTable("account_read", func(table *fazzdb.MigrationTable) {
			table.Field(fazzdb.CreateUuid("id").Primary())
			table.Field(fazzdb.CreateInteger("version"))
			table.Field(fazzdb.CreateString("name"))
			table.Field(fazzdb.CreateBigInteger("balance"))
		}),
	},
}
