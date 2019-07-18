package container

import (
	"context"
	"github.com/payfazz/go-apt/config"
	"github.com/payfazz/go-apt/example/eventsource/domain/todo"
	todoc "github.com/payfazz/go-apt/example/eventsource/domain/todo/command"
	todoq "github.com/payfazz/go-apt/example/eventsource/domain/todo/query"
	"github.com/payfazz/go-apt/pkg/fazzdb"
	"github.com/payfazz/go-apt/pkg/fazzeventsource"
	"github.com/payfazz/go-apt/pkg/fazzpubsub"
)

type ServiceContainer struct {
	TodoService todo.ServiceInterface
}

func BuildServiceContainer() *ServiceContainer {
	pubsub := fazzpubsub.NewInternalPubSub()
	return &ServiceContainer{
		TodoService: ProvideTodoService(pubsub),
	}
}

func ProvideTodoService(pubsub fazzpubsub.PubSub) todo.ServiceInterface {
	eventStore := fazzeventsource.PostgresEventStore("todo_events")
	snapshotStore := fazzeventsource.PostgresSnapshotStore("todo_snapshots")
	eventPublisher := fazzeventsource.NewEventPublisher(pubsub, "todo")
	eventRepo := todoc.NewTodoEventRepository(eventStore, snapshotStore, eventPublisher)
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
