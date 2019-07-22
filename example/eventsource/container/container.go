package container

import (
	"context"
	"github.com/payfazz/go-apt/config"
	"github.com/payfazz/go-apt/example/eventsource/domain/todo"
	todoc "github.com/payfazz/go-apt/example/eventsource/domain/todo/command"
	todoq "github.com/payfazz/go-apt/example/eventsource/domain/todo/query"
	"github.com/payfazz/go-apt/pkg/fazzdb"
	"github.com/payfazz/go-apt/pkg/fazzeventsource"
)

type ServiceContainer struct {
	TodoService todo.ServiceInterface
}

func BuildServiceContainer() *ServiceContainer {
	return &ServiceContainer{
		TodoService: ProvideTodoService(),
	}
}

func ProvideTodoService() todo.ServiceInterface {
	eventStore := fazzeventsource.PostgresEventStore("todo_events")
	snapshotStore := fazzeventsource.PostgresSnapshotStore("todo_snapshots")
	eventRepo := todoc.NewTodoEventRepository(eventStore, snapshotStore)
	command := todoc.NewTodoCommand(eventRepo)

	readModel := todoq.TodoReadModel()
	readRepo := todoq.NewTodoReadRepository(readModel)
	query := todoq.NewTodoQuery(readRepo)

	service := todo.NewTodoService(command, query)
	return service
}

func BuildContext() context.Context {
	queryDb := fazzdb.QueryDb(config.GetDB(),
		fazzdb.Config{
			Limit:           20,
			Offset:          0,
			Lock:            fazzdb.LO_NONE,
			DevelopmentMode: true,
		})

	ctx := context.Background()
	ctx = fazzdb.NewQueryContext(ctx, queryDb)

	return ctx
}
