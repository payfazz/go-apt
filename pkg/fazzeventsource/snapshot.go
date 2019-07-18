package fazzeventsource

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx/types"
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

type Snapshot struct {
	AggregateId string          `json:"aggregate_id" db:"aggregate_id"`
	LastEventId int64           `json:"last_event_id" db:"last_event_id"`
	Data        json.RawMessage `json:"data" db:"data"`
}

type SnapshotStore interface {
	Save(ctx context.Context, aggregateId string, lastEventId int64, data interface{}) (*Snapshot, error)
	FindById(ctx context.Context, aggregateId string) (*Snapshot, error)
}

type postgresSnapshotStore struct {
	tableName string
}

func (s *postgresSnapshotStore) Save(ctx context.Context, aggregateId string, lastEventId int64, data interface{}) (*Snapshot, error) {
	dataJsonByte, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	dataJsonText := types.JSONText(dataJsonByte)

	query, err := fazzdb.GetQueryContext(ctx)
	if err != nil {
		return nil, err
	}

	var ev = &Snapshot{}
	queryString := `INSERT INTO %s (aggregate_id,last_event_id,data) VALUES ($1,$2,$3) ON CONFLICT (aggregate_id) 
					DO UPDATE SET last_event_id = excluded.last_event_id, data = excluded.data RETURNING *`
	queryString = fmt.Sprintf(queryString, s.tableName)
	result, err := query.RawFirstCtx(ctx, ev, queryString, aggregateId, lastEventId, dataJsonText)
	if err != nil {
		return nil, err
	}
	return result.(*Snapshot), err
}

func (s *postgresSnapshotStore) FindById(ctx context.Context, aggregateId string) (*Snapshot, error) {
	return nil, nil
}

func PostgresSnapshotStore(tableName string) SnapshotStore {
	return &postgresSnapshotStore{tableName: tableName}
}
