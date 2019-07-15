package container

import (
	"github.com/payfazz/go-apt/example/eventsourcing/internal/domain/todo"
	todocommand "github.com/payfazz/go-apt/example/eventsourcing/internal/domain/todo/command"
	todoquery "github.com/payfazz/go-apt/example/eventsourcing/internal/domain/todo/query"
	"github.com/payfazz/go-apt/example/eventsourcing/lib/fazzeventsource"
	"github.com/payfazz/go-apt/example/eventsourcing/lib/fazzpubsub"
)

type ServiceContainer struct {
	TodoService todo.ServiceInterface
}

func BuildServiceContainer() *ServiceContainer {
	pubsub := fazzpubsub.NewInternalPubSub()
	store := fazzeventsource.NewEventStore(pubsub)
	return &ServiceContainer{
		TodoService: ProvideTodoService(store),
	}
}

func ProvideTodoService(store fazzeventsource.EventStore) todo.ServiceInterface {
	eventRepo := todocommand.NewTodoEventRepository(store)
	command := todocommand.NewTodoCommand(eventRepo)

	readModel := todoquery.TodoReadModel()
	readRepo := todoquery.NewTodoReadRepository(readModel)
	query := todoquery.NewTodoQuery(readRepo)

	service := todo.NewTodoService(command, query)
	return service
}
