package todo

import (
	_ "github.com/lib/pq"
	"github.com/payfazz/go-apt/example/eventsource/domain/todo/command"
	"github.com/payfazz/go-apt/example/eventsource/domain/todo/data"
	"github.com/payfazz/go-apt/example/eventsource/domain/todo/query"
	"github.com/payfazz/go-apt/example/eventsource/test"
	"github.com/payfazz/go-apt/pkg/fazzeventsource"
	"github.com/payfazz/go-apt/pkg/fazzpubsub"
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
		completed, text := true, "test"
		payload := data.PayloadUpdateTodo{
			Id:        *todoId,
			Completed: &completed,
			Text:      &text,
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
	eventStore := fazzeventsource.PostgresEventStore("todo_events")
	snapshotStore := fazzeventsource.PostgresSnapshotStore("todo_snapshots")
	eventPublisher := fazzeventsource.NewEventPublisher(pubsub, "todo")

	eventRepo := command.NewTodoEventRepository(eventStore, snapshotStore, eventPublisher)
	todoCommand := command.NewTodoCommand(eventRepo)

	readModel := query.TodoReadModel()
	readRepo := query.NewTodoReadRepository(readModel)
	todoQuery := query.NewTodoQuery(readRepo)

	service := NewTodoService(todoCommand, todoQuery)
	return service

}
