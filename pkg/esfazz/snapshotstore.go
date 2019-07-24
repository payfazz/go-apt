package esfazz

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx/types"
)

type Snapshot struct {
	Id      string          `json:"id" db:"id"`
	Version int             `json:"version" db:"version"`
	Data    json.RawMessage `json:"data" db:"data"`
}

type SnapshotStore interface {
	Save(ctx context.Context, data Aggregate) (*Snapshot, error)
	FindBy(ctx context.Context, id string) (*Snapshot, error)
}

type postgresSnapshotStore struct {
	tableName string
}

func (s *postgresSnapshotStore) Save(ctx context.Context, data Aggregate) (*Snapshot, error) {
	dataJsonByte, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	dataJsonText := types.JSONText(dataJsonByte)

	query, err := getContext(ctx)
	if err != nil {
		return nil, err
	}

	var ev = &Snapshot{}
	queryText := `INSERT INTO %s (id,version,data) VALUES ($1,$2,$3) ON CONFLICT (id) 
					DO UPDATE SET version = excluded.version, data = excluded.data RETURNING *`
	queryText = fmt.Sprintf(queryText, s.tableName)
	result, err := query.RawFirstCtx(ctx, ev, queryText, data.GetId(), data.GetVersion(), dataJsonText)
	if err != nil {
		return nil, err
	}
	return result.(*Snapshot), err
}

func (s *postgresSnapshotStore) FindBy(ctx context.Context, aggregateId string) (*Snapshot, error) {
	query, err := getContext(ctx)
	if err != nil {
		return nil, err
	}

	snap := &Snapshot{}
	queryText := fmt.Sprintf(`SELECT * FROM %s WHERE id = $1`, s.tableName)
	results, err := query.RawAllCtx(ctx, snap, queryText, aggregateId)
	if err != nil {
		return nil, err
	}
	snaps := results.([]*Snapshot)
	if len(snaps) == 0 {
		return nil, nil
	}
	return snaps[0], err
}

func PostgresSnapshotStore(tableName string) SnapshotStore {
	return &postgresSnapshotStore{tableName: tableName}
}
