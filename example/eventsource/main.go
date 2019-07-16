package main

import (
	"fmt"
	"github.com/payfazz/go-apt/example/eventsource/config"
	"github.com/payfazz/go-apt/example/eventsource/container"
	"github.com/payfazz/go-apt/example/eventsource/database/migration"
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

	ctn := container.BuildServiceContainer()
	ctx := container.BuildContext()

	todos, _ := ctn.TodoService.All(ctx)
	fmt.Println(todos)
}
