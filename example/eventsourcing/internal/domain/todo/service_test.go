package todo

import (
	"github.com/payfazz/go-apt/example/eventsourcing/internal/domain/todo/data"
	"github.com/payfazz/go-apt/example/eventsourcing/test"
	"testing"
)

func TestService(t *testing.T) {
	var todoId *string
	var err error
	ctx := test.PrepareTestContext()
	todoService := NewTodoService()

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

	//t.Run("All", func(t *testing.T) {
	//	result, err := todoService.All(ctx)
	//	if err != nil {
	//		t.Errorf("Error: %s", err)
	//	}
	//	if len(result) == 0 {
	//		t.Errorf("All result is empty")
	//	}
	//})

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
