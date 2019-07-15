package migration

import (
	"github.com/payfazz/go-apt/example/eventsourcing/lib/fazzeventsource"
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

var Version1 = fazzdb.MigrationVersion{
	Tables: []*fazzdb.MigrationTable{
		fazzeventsource.CreateEventsTable("events"),
	},
}
