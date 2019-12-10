package fazzeventsource

import (
	"context"
	"errors"
	"fmt"

	"github.com/payfazz/go-apt/pkg/fazzdb"
)

// RepositoryInterface is interface for event repository
type RepositoryInterface interface {
	Commit(ctx context.Context, entity EntityInterface) error
	GetLatest(ctx context.Context, entity EntityInterface) error
}

type repository struct {
	tableName string
}

// Commit will commit event to event store
func (r *repository) Commit(ctx context.Context, entity EntityInterface) error {
	query, err := fazzdb.GetTransactionOrQueryContext(ctx)
	if err != nil {
		return fmt.Errorf("es repo: %v", err)
	}

	events := entity.GetUncommittedEvents()
	baseVersion := entity.GetVersion() - len(events)

	eventLogs := make([]*EventLog, len(events))
	for i, ev := range events {
		el := EventLogModel(r.tableName)
		el.EntityId = entity.GetId()
		el.EntityVersion = baseVersion + i
		el.Event = ev
		eventLogs[i] = el
	}

	_, err = query.Use(EventLogModel(r.tableName)).BulkInsertCtx(ctx, eventLogs)
	if err != nil {
		return fmt.Errorf("es repo: %v", err)
	}

	entity.ClearUncommittedEvents()
	return nil
}

// GetLatest update entity with all event from repository
func (r *repository) GetLatest(ctx context.Context, entity EntityInterface) error {
	if len(entity.GetUncommittedEvents()) > 0 {
		return errors.New("es repo: can't get latest version, entity contain uncommitted event")
	}

	query, err := fazzdb.GetTransactionOrQueryContext(ctx)
	if err != nil {
		return fmt.Errorf("es repo: %v", err)
	}

	conditions := []fazzdb.SliceCondition{
		{Connector: fazzdb.CO_NONE, Field: "entity_id", Operator: fazzdb.OP_EQUALS, Value: entity.GetId()},
		{Connector: fazzdb.CO_AND, Field: "entity_version", Operator: fazzdb.OP_MORE_THAN_EQUALS, Value: entity.GetVersion()},
	}

	result, err := query.Use(EventLogModel(r.tableName)).
		WhereMany(conditions...).
		OrderBy("entity_version", fazzdb.DIR_ASC).
		AllCtx(ctx)
	if err != nil {
		return fmt.Errorf("es repo: %v", err)
	}

	eventLogs := result.([]*EventLog)
	for _, el := range eventLogs {
		err = entity.Apply(el.Event)
		if err != nil {
			return fmt.Errorf("es repo: %v", err)
		}
	}

	entity.ClearUncommittedEvents()
	return nil
}

// NewRepository return new event repository
func NewRepository(tableName string) RepositoryInterface {
	return &repository{tableName: tableName}
}
