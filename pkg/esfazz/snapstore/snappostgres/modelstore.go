package snappostgres

import (
	"context"
	"encoding/json"
	"github.com/payfazz/go-apt/pkg/esfazz"
	"github.com/payfazz/go-apt/pkg/esfazz/snapstore"
	"github.com/payfazz/go-apt/pkg/fazzdb"
	"reflect"
)

type modelSnapshotStore struct {
	tableName    string
	model        esfazz.AggregateModel
	newAggregate esfazz.AggregateFactory
}

// Save is a function to save aggregate to database
func (s *modelSnapshotStore) Save(ctx context.Context, id string, data json.RawMessage) error {
	query, err := fazzdb.GetTransactionOrQueryContext(ctx)
	if err != nil {
		return err
	}

	agg := s.newAggregate(id)
	err = json.Unmarshal(data, agg)
	if err != nil {
		return err
	}

	count, err := query.Use(s.model).
		Where("id", id).
		WithLimit(0).
		Count()
	if err != nil {
		return err
	}

	if *count == 0 {
		_, err = query.Use(agg.(fazzdb.ModelInterface)).InsertCtx(ctx, false)
	} else {
		_, err = query.Use(agg.(fazzdb.ModelInterface)).UpdateCtx(ctx)
	}
	return err
}

// Find find aggregate in database based on id
func (s *modelSnapshotStore) Find(ctx context.Context, id string) (json.RawMessage, error) {
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

	items := reflect.ValueOf(row)
	if items.Len() == 0 {
		return nil, nil
	}

	result, err := json.Marshal(items.Index(0).Interface())
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ModelSnapshotStore is a constructor for PostgreSQL Model based snapshot store
func ModelSnapshotStore(newAggregate esfazz.AggregateFactory) snapstore.SnapshotStore {
	return &modelSnapshotStore{
		newAggregate: newAggregate,
		model:        newAggregate("").(esfazz.AggregateModel),
	}
}
