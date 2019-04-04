package migration

import "github.com/payfazz/go-apt/pkg/fazzdb"

var Version1 = fazzdb.MigrationVersion{
	Tables: []*fazzdb.MigrationTable{
		fazzdb.CreateTable("students", func(table *fazzdb.MigrationTable) {
			table.Field(fazzdb.CreateSerial("id").Primary())
			table.Field(fazzdb.CreateString("name"))
			table.Field(fazzdb.CreateInteger("number"))
			table.Timestamps()
			table.SoftDelete()
		}),
		fazzdb.CreateTable("classes", func(table *fazzdb.MigrationTable) {
			table.Field(fazzdb.CreateSerial("id").Primary())
			table.Field(fazzdb.CreateString("name"))
			table.Timestamps()
			table.SoftDelete()
		}),
		fazzdb.CreateTable("student_classes", func(table *fazzdb.MigrationTable) {
			table.Field(fazzdb.CreateInteger("student_id").Primary())
			table.Field(fazzdb.CreateInteger("class_id").Primary())

			table.Reference(
				fazzdb.Foreign("student_id").
					Reference("id").
					On("students").
					OnDelete(fazzdb.RA_CASCADE).
					OnUpdate(fazzdb.RA_CASCADE),
			)
			table.Reference(
				fazzdb.Foreign("class_id").
					Reference("id").
					On("classes").
					OnDelete(fazzdb.RA_CASCADE).
					OnUpdate(fazzdb.RA_CASCADE),
			)
		}),
	},
}
