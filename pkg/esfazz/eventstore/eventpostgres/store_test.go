package eventpostgres

import (
	"context"
	_ "github.com/lib/pq"
	"github.com/payfazz/go-apt/config"
	"github.com/payfazz/go-apt/pkg/esfazz"
	"github.com/payfazz/go-apt/pkg/fazzdb"
	"testing"
)

func TestPostgresEventStore_Save(t *testing.T) {
	ctx := prepareContext()

	store := EventStore("event")

	agg := &esfazz.BaseAggregate{
		Id:      "01234567-89ab-cdef-0123-456789abcdef",
		Version: 0,
	}

	evs, err := store.Save(ctx, agg, &esfazz.EventPayload{
		Type: "test.event",
		Data: map[string]string{"test": "234"},
	})

	if err != nil {
		t.Errorf("error saving data: %s", err)
	}

	if len(evs) == 0 {
		t.Error("save result empty")
	}
}

func TestPostgresEventStore_FindNotApplied(t *testing.T) {
	ctx := prepareContext()

	store := EventStore("event")

	agg := &esfazz.BaseAggregate{
		Id:      "01234567-89ab-cdef-0123-456789abcdef",
		Version: 0,
	}

	events := []*esfazz.EventPayload{
		{
			Type: "test.event",
			Data: map[string]string{"test": "234"},
		},
		{
			Type: "test.event",
			Data: map[string]string{"test": "345"},
		},
		{
			Type: "test.event",
			Data: map[string]string{"test": "456"},
		},
	}

	evs, err := store.Save(ctx, agg, events...)
	if err != nil {
		t.Errorf("error saving data: %s", err)
	}
	if len(evs) != 3 {
		t.Errorf("saved event list not of same length, expected: 3, result %d", len(evs))
	}

	evResults, err := store.FindNotApplied(ctx, agg)

	if err != nil {
		t.Errorf("error saving data: %s", err)
	}
	if len(evResults) != 3 {
		t.Errorf("data in list not on the same length as existing data, expected: 3, result: %d", len(evResults))
	}

	for idx := range events {
		if evResults[idx].Aggregate.GetVersion() != int64(idx+1) {
			t.Errorf("event in index %d is not in expected order", idx)
		}
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
