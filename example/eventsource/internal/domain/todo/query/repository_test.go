package query

import (
	"github.com/gofrs/uuid"
	"github.com/payfazz/go-apt/example/eventsource/test"
	"testing"
)

func TestTodoReadRepository(t *testing.T) {
	var err error
	ctx := test.PrepareTestContext()
	repo := NewTodoReadRepository(TodoReadModel())

	uuidV4, _ := uuid.NewV4()
	id := uuidV4.String()

	t.Run("Create", func(t *testing.T) {
		model := TodoReadModel()
		model.Id = id
		model.Text = "test"
		model.Completed = false

		resultId, err := repo.Create(ctx, model)
		if err != nil {
			t.Errorf("Create error: %s", err)
		}
		if resultId == nil {
			t.Errorf("Create result is empty")
		}
	})

	t.Run("All", func(t *testing.T) {
		result, err := repo.All(ctx)
		if err != nil {
			t.Errorf("All error: %s", err)
		}
		if len(result) == 0 {
			t.Errorf("All result is empty")
		}
	})

	t.Run("Update", func(t *testing.T) {
		model := TodoReadModel()
		model.Id = id
		model.Text = "test"
		model.Completed = false

		err = repo.Update(ctx, model)
		if err != nil {
			t.Errorf("Update error: %s", err)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		model := TodoReadModel()
		model.Id = id

		err = repo.Delete(ctx, model)
		if err != nil {
			t.Errorf("Delete error: %s", err)
		}
	})

}
