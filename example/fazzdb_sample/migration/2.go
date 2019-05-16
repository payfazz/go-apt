package migration

import (
	"github.com/payfazz/go-apt/example/fazzdb_sample/seed"
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

var bookStatusEnum = fazzdb.NewEnum(
	"book_status",
	"BORROWED", "AVAILABLE",
)

var Version2 = fazzdb.MigrationVersion{
	Enums: []*fazzdb.MigrationEnum{
		bookStatusEnum,
	},
	Tables: []*fazzdb.MigrationTable{
		fazzdb.AlterTable("books", func(table *fazzdb.MigrationTable) {
			table.Field(fazzdb.AddEnum("status", bookStatusEnum))
			table.Field(fazzdb.AddInteger("year"))
		}),
	},
	Seeds: []fazzdb.SeederInterface{
		seed.AuthorSeeder(),
		seed.BookSeeder(),
	},
	Raw: `CREATE TABLE raw_queries (id serial primary key, name varchar);`,
}
