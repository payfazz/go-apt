package fazzeventsource

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
	Save(ctx context.Context, event *EventLog, data interface{}) (*Snapshot, error)
	FindBy(ctx context.Context, aggregateId string) (*Snapshot, error)
}

type postgresSnapshotStore struct {
	tableName string
}

func (s *postgresSnapshotStore) Save(ctx context.Context, event *EventLog, data interface{}) (*Snapshot, error) {
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
	queryString := `INSERT INTO %s (id,version,data) VALUES ($1,$2,$3) ON CONFLICT (id) 
					DO UPDATE SET version = excluded.version, data = excluded.data RETURNING *`
	queryString = fmt.Sprintf(queryString, s.tableName)
	result, err := query.RawFirstCtx(ctx, ev, queryString, event.AggregateId, event.AggregateVersion, dataJsonText)
	if err != nil {
		return nil, err
	}
	return result.(*Snapshot), err
}

func (s *postgresSnapshotStore) FindBy(ctx context.Context, aggregateId string) (*Snapshot, error) {
	return nil, nil
}

func PostgresSnapshotStore(tableName string) SnapshotStore {
	return &postgresSnapshotStore{tableName: tableName}
}
