package fazzeventsource

import (
	"github.com/gofrs/uuid"
	"github.com/payfazz/go-apt/example/fazzeventsource_sample/test"
	"testing"
)

func TestEventRepository_Save(t *testing.T) {
	ctx := test.PrepareTestContext()
	store := NewPostgresEventStore("events")

	_, err := store.Save(ctx, "test.event", map[string]string{"id": "123", "test": "234"})
	if err != nil {
		t.Errorf("error saving data: %s", err)
	}
}

func TestEventRepository_AggregateById(t *testing.T) {
	ctx := test.PrepareTestContext()
	store := NewPostgresEventStore("events")

	uuidV4, _ := uuid.NewV4()
	id := uuidV4.String()

	_, _ = store.Save(ctx, "test.event", map[string]string{"id": id, "test": "234"})
	_, _ = store.Save(ctx, "test.event", map[string]string{"id": id, "test": "345"})
	_, _ = store.Save(ctx, "test.event", map[string]string{"id": id, "test": "456"})

	evs, err := store.FindByInstanceId(ctx, id)
	if err != nil {
		t.Errorf("error getting data: %s", err)
	}

	if len(evs) != 3 {
		t.Errorf("aggregate data length not valid, should be 3 not %d", len(evs))
	}
}
