package esfazz

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx/types"
	"time"
)

// EventStore is an interface used for event store
type EventStore interface {
	Save(ctx context.Context, ev EventPayload) (*EventLog, error)
	FindAllBy(ctx context.Context, aggregateId string, firstVersion int) ([]*EventLog, error)
}

type postgresEventStore struct {
	tableName string
}

// Save is a function to save event to event store
func (e *postgresEventStore) Save(ctx context.Context, ev EventPayload) (*EventLog, error) {

	// if no aggregate, event will be related to new aggregate object
	if ev.Aggregate == nil {
		uuidV4, _ := uuid.NewV4()
		ev.Aggregate = &BaseAggregate{
			Id:      uuidV4.String(),
			Version: 0,
		}
	}

	latestVersion, err := e.findLatestVersion(ctx, ev.Aggregate.GetId())
	if err != nil {
		return nil, err
	}
	if latestVersion >= ev.Aggregate.GetVersion() {
		return nil, errors.New("event aggregate version is lower than latest save event")
	}

	dataJsonByte, err := json.Marshal(ev.Data)
	if err != nil {
		return nil, err
	}
	dataJsonText := types.JSONText(dataJsonByte)

	query, err := getContext(ctx)
	if err != nil {
		return nil, err
	}

	el := &EventLog{}
	queryText := fmt.Sprintf(`INSERT INTO %s (event_type, aggregate_id, aggregate_version, data, created_at) 
									VALUES ($1,$2,$3,$4,$5) RETURNING *`, e.tableName)
	result, err := query.RawFirstCtx(ctx, el, queryText, ev.Type,
		ev.Aggregate.GetId(), ev.Aggregate.GetVersion(), dataJsonText, time.Now())
	if err != nil {
		return nil, err
	}
	return result.(*EventLog), err
}

// FindAllBy return all event filtered by aggregateId and version
func (e *postgresEventStore) FindAllBy(ctx context.Context, aggregateId string, firstVersion int) ([]*EventLog, error) {
	query, err := getContext(ctx)
	if err != nil {
		return nil, err
	}
	el := &EventLog{}
	queryText := fmt.Sprintf(`SELECT * FROM %s WHERE aggregate_id = $1 AND aggregate_version >= $2 
									ORDER BY event_id ASC`, e.tableName)
	results, err := query.RawAllCtx(ctx, el, queryText, aggregateId, firstVersion)
	if err != nil {
		return nil, err
	}
	return results.([]*EventLog), err
}

func (e *postgresEventStore) findLatestVersion(ctx context.Context, aggregateId string) (int, error) {
	query, err := getContext(ctx)
	if err != nil {
		return 0, err
	}

	el := &EventLog{}
	queryText := fmt.Sprintf(`SELECT * FROM %s WHERE aggregate_id = $1 ORDER BY event_id DESC LIMIT 1`, e.tableName)
	results, err := query.RawAllCtx(ctx, el, queryText, aggregateId)
	if err != nil {
		return 0, err
	}

	els := results.([]*EventLog)
	if len(els) == 0 {
		return -1, nil
	}
	return els[0].AggregateVersion, err

}

// PostgresEventStore is a function to create new EventStore
func PostgresEventStore(tableName string) EventStore {
	return &postgresEventStore{
		tableName: tableName,
	}
}
