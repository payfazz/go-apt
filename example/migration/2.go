package migration

import "github.com/payfazz/go-apt/pkg/fazzdb"

var attendanceEnum = fazzdb.NewEnum(
	"attendance",
	"PRESENT",
	"NOT PRESENT",
)

var Version2 = fazzdb.MigrationVersion{
	Enums: []*fazzdb.MigrationEnum{
		attendanceEnum,
	},
	Tables: []*fazzdb.MigrationTable{
		fazzdb.AlterTable("students", func(table *fazzdb.MigrationTable) {
			table.Field(fazzdb.AddEnum("attendances", *attendanceEnum))
			table.Field(fazzdb.AlterDecimal("number", 20, 2))
			table.Field(fazzdb.DropColumn("createdAt"))
		}),
		fazzdb.AlterTable("students", func(table *fazzdb.MigrationTable) {
			table.Field(fazzdb.RenameColumn("name", "studentName"))
		}),
	},
}
