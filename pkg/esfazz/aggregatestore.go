package esfazz

import (
	"context"
	"encoding/json"
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

// AggregateStore is interface for aggregate storage
type AggregateStore interface {
	Save(ctx context.Context, agg Aggregate) (*AggregateRow, error)
	Find(ctx context.Context, id string) (*AggregateRow, error)
}

type postgresAggregateStore struct {
	tableName string
	model     *AggregateRow
}

// Save is a function to save aggregate to database
func (s *postgresAggregateStore) Save(ctx context.Context, agg Aggregate) (*AggregateRow, error) {
	query, err := fazzdb.GetTransactionOrQueryContext(ctx)
	if err != nil {
		return nil, err
	}

	dataJsonByte, err := json.Marshal(agg)
	if err != nil {
		return nil, err
	}
	data := json.RawMessage(dataJsonByte)

	updateRow := AggregateRowModel(s.tableName)
	updateRow.Id = agg.GetId()
	updateRow.Version = agg.GetVersion()
	updateRow.Data = data

	count, err := query.Use(s.model).
		Where("id", agg.GetId()).
		WithLimit(0).
		Count()
	if err != nil {
		return nil, err
	}

	if *count == 0 {
		_, err = query.Use(updateRow).InsertCtx(ctx, false)
	} else {
		_, err = query.Use(updateRow).UpdateCtx(ctx)
	}
	if err != nil {
		return nil, err
	}

	return updateRow, err
}

// Find find aggregate in database based on id
func (s *postgresAggregateStore) Find(ctx context.Context, id string) (*AggregateRow, error) {
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

	return results[0], nil
}

// PostgresAggregateStore is a constructor for PostgreSQL based aggregate store
func PostgresAggregateStore(tableName string) AggregateStore {
	return &postgresAggregateStore{
		tableName: tableName,
		model:     AggregateRowModel(tableName),
	}
}
