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

	evs, err := store.Save(ctx, &esfazz.EventPayload{
		Type: "test.event",
		Aggregate: &esfazz.BaseAggregate{
			Id:      "01234567-89ab-cdef-0123-456789abcdef",
			Version: 0,
		},
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
	events := []*esfazz.EventPayload{
		{
			Type: "test.event",
			Aggregate: &esfazz.BaseAggregate{
				Id:      "01234567-89ab-cdef-0123-456789abcdef",
				Version: 0,
			},
			Data: map[string]string{"test": "234"},
		},
		{
			Type: "test.event",
			Aggregate: &esfazz.BaseAggregate{
				Id:      "01234567-89ab-cdef-0123-456789abcdef",
				Version: 1,
			},
			Data: map[string]string{"test": "345"},
		},
		{
			Type: "test.event",
			Aggregate: &esfazz.BaseAggregate{
				Id:      "01234567-89ab-cdef-0123-456789abcdef",
				Version: 2,
			},
			Data: map[string]string{"test": "456"},
		},
	}

	evs, err := store.Save(ctx, events...)
	if err != nil {
		t.Errorf("error saving data: %s", err)
	}
	if len(evs) != 3 {
		t.Errorf("saved event list not of same length, expected: 3, result %d", len(evs))
	}

	evResults, err := store.FindNotApplied(ctx,
		&esfazz.BaseAggregate{
			Id:      "01234567-89ab-cdef-0123-456789abcdef",
			Version: 0,
		},
	)

	if err != nil {
		t.Errorf("error saving data: %s", err)
	}
	if len(evResults) != 3 {
		t.Errorf("data in list not on the same length as existing data, expected: 3, result: %d", len(evResults))
	}

	for idx := range events {
		if events[idx].Aggregate.GetVersion() != evResults[idx].Aggregate.GetVersion() {
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
