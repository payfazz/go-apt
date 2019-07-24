package esfazz

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx/types"
)

type AggregateRow struct {
	Id      string          `json:"id" db:"id"`
	Version int             `json:"version" db:"version"`
	Data    json.RawMessage `json:"data" db:"data"`
}

type AggregateStore interface {
	Save(ctx context.Context, data Aggregate) (*AggregateRow, error)
	FindBy(ctx context.Context, id string) (*AggregateRow, error)
}

type postgresAggregateStore struct {
	tableName string
}

func (s *postgresAggregateStore) Save(ctx context.Context, data Aggregate) (*AggregateRow, error) {
	dataJsonByte, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	dataJsonText := types.JSONText(dataJsonByte)

	query, err := getContext(ctx)
	if err != nil {
		return nil, err
	}

	var ev = &AggregateRow{}
	queryText := `INSERT INTO %s (id,version,data) VALUES ($1,$2,$3) ON CONFLICT (id) 
					DO UPDATE SET version = excluded.version, data = excluded.data RETURNING *`
	queryText = fmt.Sprintf(queryText, s.tableName)
	result, err := query.RawFirstCtx(ctx, ev, queryText, data.GetId(), data.GetVersion(), dataJsonText)
	if err != nil {
		return nil, err
	}
	return result.(*AggregateRow), err
}

func (s *postgresAggregateStore) FindBy(ctx context.Context, aggregateId string) (*AggregateRow, error) {
	query, err := getContext(ctx)
	if err != nil {
		return nil, err
	}

	snap := &AggregateRow{}
	queryText := fmt.Sprintf(`SELECT * FROM %s WHERE id = $1`, s.tableName)
	results, err := query.RawAllCtx(ctx, snap, queryText, aggregateId)
	if err != nil {
		return nil, err
	}
	snaps := results.([]*AggregateRow)
	if len(snaps) == 0 {
		return nil, nil
	}
	return snaps[0], err
}

func PostgresAggregateStore(tableName string) AggregateStore {
	return &postgresAggregateStore{tableName: tableName}
}
