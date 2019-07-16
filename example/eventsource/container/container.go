package container

import (
	"context"
	"github.com/payfazz/go-apt/example/eventsource/config"
	"github.com/payfazz/go-apt/example/eventsource/internal/domain/todo"
	todocommand "github.com/payfazz/go-apt/example/eventsource/internal/domain/todo/command"
	todoquery "github.com/payfazz/go-apt/example/eventsource/internal/domain/todo/query"
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
	eventStore := fazzeventsource.NewPostgresEventStore("events")
	eventPublisher := fazzeventsource.NewEventPublisher(pubsub, "todo")
	eventRepo := todocommand.NewTodoEventRepository(eventStore, eventPublisher)
	command := todocommand.NewTodoCommand(eventRepo)

	readModel := todoquery.TodoReadModel()
	readRepo := todoquery.NewTodoReadRepository(readModel)
	query := todoquery.NewTodoQuery(readRepo)

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
