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

	_, err := store.Save(ctx, EventPayload{
		Type:             "test.event",
		AggregateVersion: 0,
		AggregateId:      "123",
		Data:             map[string]interface{}{"test": "234"},
	})
	if err != nil {
		t.Errorf("error saving data: %s", err)
	}
}

func TestEventStore_FindByKey(t *testing.T) {
	ctx := test.PrepareTestContext()
	store := PostgresEventStore("todo_events")

	uuidV4, _ := uuid.NewV4()
	id := uuidV4.String()

	_, _ = store.Save(ctx, EventPayload{
		Type:             "test.event",
		AggregateVersion: 0,
		AggregateId:      id,
		Data:             map[string]interface{}{"test": "234"},
	})
	_, _ = store.Save(ctx, EventPayload{
		Type:             "test.event",
		AggregateVersion: 1,
		AggregateId:      id,
		Data:             map[string]interface{}{"test": "345"},
	})
	_, _ = store.Save(ctx, EventPayload{
		Type:             "test.event",
		AggregateVersion: 2,
		AggregateId:      id,
		Data:             map[string]interface{}{"test": "456"},
	})

	evs, err := store.FindAllBy(ctx, id, 0)
	if err != nil {
		t.Errorf("error getting data: %s", err)
	}

	if len(evs) != 3 {
		t.Errorf("aggregate data length not valid, should be 3 not %d", len(evs))
	}
}
