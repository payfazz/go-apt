package migration

import "github.com/payfazz/go-apt/pkg/fazzdb"

var Version1 = fazzdb.MigrationVersion{
	Tables: []*fazzdb.MigrationTable{
		fazzdb.CreateTable("authors", func(table *fazzdb.MigrationTable) {
			table.Field(fazzdb.CreateSerial("id").Primary())
			table.Field(fazzdb.CreateString("name", 50))
			table.Field(fazzdb.CreateString("country", 100))
			table.Timestamps()
			table.SoftDelete()
			table.Indexes("name", "country")
		}),
		fazzdb.CreateTable("books", func(table *fazzdb.MigrationTable) {
			table.Field(fazzdb.CreateString("id", 60).Primary())
			table.Field(fazzdb.CreateString("title", 30))
			table.Field(fazzdb.CreateInteger("stock"))
			table.Field(fazzdb.CreateInteger("authorId"))
			table.Timestamps()

			table.Reference(
				fazzdb.Foreign("authorId").
					Reference("id").
					On("authors").
					OnDelete(fazzdb.RA_CASCADE).
					OnUpdate(fazzdb.RA_CASCADE),
			)
		}),
	},
}
