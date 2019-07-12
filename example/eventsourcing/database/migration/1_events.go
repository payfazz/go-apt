package migration

import "github.com/payfazz/go-apt/pkg/fazzdb"

var Version1 = fazzdb.MigrationVersion{
	Tables: []*fazzdb.MigrationTable{
		fazzdb.CreateTable("events", func(table *fazzdb.MigrationTable) {
			table.Field(fazzdb.CreateBigSerial("id").Primary())
			table.Field(fazzdb.CreateString("type"))
			table.Field(fazzdb.CreateJsonB("data"))
			table.Field(fazzdb.CreateTimestampTz("created_at", 0))
		}),
	},
}
