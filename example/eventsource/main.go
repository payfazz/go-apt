package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/payfazz/go-apt/config"
	"github.com/payfazz/go-apt/example/eventsource/container"
	"github.com/payfazz/go-apt/example/eventsource/database/migration"
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

func main() {
	fazzdb.Migrate(config.GetDB(), "es-example", true, true,
		migration.Version1,
	)

	ctn := container.BuildServiceContainer()
	ctx := container.BuildContext()

	todos, _ := ctn.TodoService.All(ctx)
	fmt.Println(todos)
}
