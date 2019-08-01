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

	err := store.Save(ctx, &esfazz.Event{
		Type: "test.event",
		Aggregate: &esfazz.BaseAggregate{
			Id:      "01234567-89ab-cdef-0123-456789abcdef",
			Version: 0,
		},
		Data: []byte("{\"test\": \"234\"}"),
	})

	if err != nil {
		t.Errorf("error saving data: %s", err)
	}
}

func TestPostgresEventStore_FindNotApplied(t *testing.T) {
	ctx := prepareContext()

	store := EventStore("event")
	err := store.Save(ctx, &esfazz.Event{
		Type: "test.event",
		Aggregate: &esfazz.BaseAggregate{
			Id:      "01234567-89ab-cdef-0123-456789abcdef",
			Version: 0,
		},
		Data: []byte("{\"test\": \"234\"}"),
	})
	if err != nil {
		t.Errorf("error saving data: %s", err)
	}

	err = store.Save(ctx, &esfazz.Event{
		Type: "test.event",
		Aggregate: &esfazz.BaseAggregate{
			Id:      "01234567-89ab-cdef-0123-456789abcdef",
			Version: 1,
		},
		Data: []byte("{\"test\": \"345\"}"),
	})
	if err != nil {
		t.Errorf("error saving data: %s", err)
	}

	err = store.Save(ctx, &esfazz.Event{
		Type: "test.event",
		Aggregate: &esfazz.BaseAggregate{
			Id:      "01234567-89ab-cdef-0123-456789abcdef",
			Version: 2,
		},
		Data: []byte("{\"test\": \"456\"}"),
	})
	if err != nil {
		t.Errorf("error saving data: %s", err)
	}

	evs, err := store.FindNotApplied(ctx,
		&esfazz.BaseAggregate{
			Id:      "01234567-89ab-cdef-0123-456789abcdef",
			Version: 0,
		},
	)
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
