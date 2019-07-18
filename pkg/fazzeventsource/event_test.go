package fazzeventsource

import (
	"github.com/gofrs/uuid"
	_ "github.com/lib/pq"
	"github.com/payfazz/go-apt/example/eventsource/test"
	"testing"
)

func TestEventStore_Save(t *testing.T) {
	ctx := test.PrepareTestContext()
	store := PostgresEventStore("todo_events")

	_, err := store.Save(ctx, "test.event", map[string]string{"id": "123", "test": "234"})
	if err != nil {
		t.Errorf("error saving data: %s", err)
	}
}

func TestEventStore_FindByKey(t *testing.T) {
	ctx := test.PrepareTestContext()
	store := PostgresEventStore("todo_events")

	uuidV4, _ := uuid.NewV4()
	id := uuidV4.String()

	_, _ = store.Save(ctx, "test.event", map[string]string{"id": id, "test": "234"})
	_, _ = store.Save(ctx, "test.event", map[string]string{"id": id, "test": "345"})
	_, _ = store.Save(ctx, "test.event", map[string]string{"id": id, "test": "456"})

	evs, err := store.FindAllByKey(ctx, "id", id, 0)
	if err != nil {
		t.Errorf("error getting data: %s", err)
	}

	if len(evs) != 3 {
		t.Errorf("aggregate data length not valid, should be 3 not %d", len(evs))
	}
}
