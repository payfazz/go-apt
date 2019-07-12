package main

import (
	"github.com/payfazz/go-apt/example/eventsourcing/config"
	"github.com/payfazz/go-apt/example/eventsourcing/database/migration"
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

func main() {
	fazzdb.Migrate(config.GetDB(),
		"estodo-example",
		true,
		true,
		migration.Version1,
		migration.Version2,
	)
}
