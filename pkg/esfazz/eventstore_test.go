package esfazz

import (
	"context"
	_ "github.com/lib/pq"
	"github.com/payfazz/go-apt/config"
	"github.com/payfazz/go-apt/pkg/fazzdb"
	"testing"
)

func TestEventStore_Save(t *testing.T) {
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

	store := PostgresEventStore("event")
	_, err := store.Save(ctx, EventPayload{
		Type: "test.event",
		Data: map[string]interface{}{"test": "234"},
	})
	if err != nil {
		t.Errorf("error saving data: %s", err)
	}
}
