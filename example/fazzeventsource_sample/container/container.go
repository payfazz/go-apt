package container

import (
	"github.com/payfazz/go-apt/example/fazzeventsource_sample/internal/domain/todo"
	todocommand "github.com/payfazz/go-apt/example/fazzeventsource_sample/internal/domain/todo/command"
	todoquery "github.com/payfazz/go-apt/example/fazzeventsource_sample/internal/domain/todo/query"
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
