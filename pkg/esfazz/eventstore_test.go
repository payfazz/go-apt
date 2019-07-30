package esfazz

import (
	"context"
	_ "github.com/lib/pq"
	"github.com/payfazz/go-apt/config"
	"github.com/payfazz/go-apt/pkg/fazzdb"
	"testing"
)

func TestPostgresEventStore_Save(t *testing.T) {
	ctx := prepareContext()

	store := PostgresEventStore("event")
	_, err := store.Save(ctx, &EventPayload{
		Type: "test.event",
		Data: map[string]interface{}{"test": "234"},
	})
	if err != nil {
		t.Errorf("error saving data: %s", err)
	}
}

func TestPostgresEventStore_FindAfterAggregate(t *testing.T) {
	ctx := prepareContext()

	store := PostgresEventStore("event")
	ev, err := store.Save(ctx, &EventPayload{
		Type: "test.event",
		Data: map[string]interface{}{"test": "234"},
	})
	if err != nil {
		t.Errorf("error saving data: %s", err)
	}

	ev, err = store.Save(ctx, &EventPayload{
		Type: "test.event",
		Aggregate: &BaseAggregate{
			Id:      ev.AggregateId,
			Version: ev.AggregateVersion + 1,
		},
		Data: map[string]interface{}{"test": "345"},
	})
	if err != nil {
		t.Errorf("error saving data: %s", err)
	}

	ev, err = store.Save(ctx, &EventPayload{
		Type: "test.event",
		Aggregate: &BaseAggregate{
			Id:      ev.AggregateId,
			Version: ev.AggregateVersion + 1,
		},
		Data: map[string]interface{}{"test": "456"},
	})
	if err != nil {
		t.Errorf("error saving data: %s", err)
	}

	evs, err := store.FindAfterAggregate(ctx, &BaseAggregate{
		Id:      ev.AggregateId,
		Version: 0,
	})
	if err != nil {
		t.Errorf("error saving data: %s", err)
	}
	if len(evs) != 3 {
		t.Errorf("existing data not list not on the same length, expected: 3, result: %d", len(evs))
	}
}

func prepareContext() context.Context {
	fazzdb.Migrate(
		config.GetDB(),
		"test-esfazz",
		true,
		true,
		fazzdb.MigrationVersion{
			Tables: []*fazzdb.MigrationTable{
				CreateEventsTable("event"),
			},
		},
	)
	queryDb := fazzdb.QueryDb(config.GetDB(), config.Parameter)
	ctx := context.Background()
	ctx = fazzdb.NewQueryContext(ctx, queryDb)
	return ctx
}
