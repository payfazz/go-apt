package snappostgres

import (
	"context"
	"encoding/json"
	"github.com/jmoiron/sqlx/types"
	"github.com/payfazz/go-apt/pkg/esfazz/snapstore"
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

type jsonSnapshotStore struct {
	tableName string
	model     *AggregateRow
}

// Save is a function to save aggregate to database
func (s *jsonSnapshotStore) Save(ctx context.Context, id string, data json.RawMessage) error {
	query, err := fazzdb.GetTransactionOrQueryContext(ctx)
	if err != nil {
		return err
	}
	dataJson := types.JSONText(data)

	updateRow := AggregateRowModel(s.tableName)
	updateRow.Id = id
	updateRow.Data = dataJson

	count, err := query.Use(s.model).
		Where("id", id).
		WithLimit(0).
		Count()
	if err != nil {
		return err
	}

	if *count == 0 {
		_, err = query.Use(updateRow).InsertCtx(ctx, false)
	} else {
		_, err = query.Use(updateRow).UpdateCtx(ctx)
	}
	return err
}

// Find find aggregate in database based on id
func (s *jsonSnapshotStore) Find(ctx context.Context, id string) (json.RawMessage, error) {
	query, err := fazzdb.GetTransactionOrQueryContext(ctx)
	if err != nil {
		return nil, err
	}
	row, err := query.Use(s.model).
		Where("id", id).
		WithLimit(1).
		AllCtx(ctx)

	if err != nil {
		return nil, err
	}

	results := row.([]*AggregateRow)
	if len(results) == 0 {
		return nil, nil
	}

	return json.RawMessage(results[0].Data), nil
}

// JSONSnapshotStore is a constructor for PostgreSQL JSON based snapshot store
func JSONSnapshotStore(tableName string) snapstore.SnapshotStore {
	return &jsonSnapshotStore{
		tableName: tableName,
		model:     AggregateRowModel(tableName),
	}
}
