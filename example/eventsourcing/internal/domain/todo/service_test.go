package todo

import (
	"github.com/payfazz/go-apt/example/eventsourcing/internal/domain/todo/command"
	"github.com/payfazz/go-apt/example/eventsourcing/internal/domain/todo/data"
	"github.com/payfazz/go-apt/example/eventsourcing/internal/domain/todo/query"
	"github.com/payfazz/go-apt/example/eventsourcing/lib/fazzeventsource"
	"github.com/payfazz/go-apt/example/eventsourcing/lib/fazzpubsub"
	"github.com/payfazz/go-apt/example/eventsourcing/test"
	"testing"
)

func TestService(t *testing.T) {
	var todoId *string
	var err error
	ctx := test.PrepareTestContext()
	todoService := provideTodoService()

	t.Run("Create", func(t *testing.T) {
		payload := data.PayloadCreateTodo{Text: "test"}
		todoId, err = todoService.Create(ctx, payload)
		if err != nil {
			t.Errorf("Error: %s", err)
		}
		if todoId == nil {
			t.Errorf("Create result is empty")
		}
	})

	t.Run("All", func(t *testing.T) {
		_, err := todoService.All(ctx)
		if err != nil {
			t.Errorf("Error: %s", err)
		}
	})

	t.Run("Update", func(t *testing.T) {
		payload := data.PayloadUpdateTodo{
			Id:        *todoId,
			Completed: true,
			Text:      "test",
		}
		err = todoService.Update(ctx, payload)
		if err != nil {
			t.Errorf("Error: %s", err)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		err = todoService.Delete(ctx, *todoId)
		if err != nil {
			t.Errorf("Error: %s", err)
		}
	})

}

func provideTodoService() ServiceInterface {
	pubsub := fazzpubsub.NewInternalPubSub()
	store := fazzeventsource.NewEventStore(pubsub)

	eventRepo := command.NewTodoEventRepository(store)
	todoCommand := command.NewTodoCommand(eventRepo)

	readModel := query.TodoReadModel()
	readRepo := query.NewTodoReadRepository(readModel)
	todoQuery := query.NewTodoQuery(readRepo)

	service := NewTodoService(todoCommand, todoQuery)
	return service

}
