package main

import (
	"github.com/payfazz/go-apt/example/fazzeventsource_sample/config"
	"github.com/payfazz/go-apt/example/fazzeventsource_sample/database/migration"
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

func main() {
	fazzdb.Migrate(config.GetDB(),
		"esexample",
		true,
		true,
		migration.Version1,
		migration.Version2,
	)
}
