package migration

import "github.com/payfazz/go-apt/pkg/fazzdb"

var Version2 = fazzdb.MigrationVersion{
	Tables: []*fazzdb.MigrationTable{
		fazzdb.CreateTable("todos_read", func(table *fazzdb.MigrationTable) {
			table.Field(fazzdb.CreateUuid("id").Primary())
			table.Field(fazzdb.CreateString("text"))
			table.Field(fazzdb.CreateBoolean("completed"))
		}),
		fazzdb.CreateTable("read_metas", func(table *fazzdb.MigrationTable) {
			table.Field(fazzdb.CreateText("key").Primary())
			table.Field(fazzdb.CreateText("value"))
		}),
	},
}
