package fazzdb

import (
	"github.com/payfazz/go-apt/pkg/fazzdb/fazzorder"
	"github.com/payfazz/go-apt/pkg/fazzdb/fazzspec"
)

func (q *Query) Where(key string, value interface{}) *Query {
	return q.AppendCondition(fazzspec.CO_AND, key, fazzspec.OP_EQUALS, value)
}

func (q *Query) WhereOp(key string, operator fazzspec.Operator, value interface{}) *Query {
	return q.AppendCondition(fazzspec.CO_AND, key, operator, value)
}

func (q *Query) OrWhere(key string, value interface{}) *Query {
	return q.AppendCondition(fazzspec.CO_OR, key, fazzspec.OP_EQUALS, value)
}

func (q *Query) OrWhereOp(key string, operator fazzspec.Operator, value interface{}) *Query {
	return q.AppendCondition(fazzspec.CO_OR, key, operator, value)
}

func (q *Query) GroupWhere(conditionFunc func(query *Query) *Query) *Query {
	query := QueryTx(q.Tx).Use(q.Model)
	param := conditionFunc(query).Parameter
	q.appendGroupConditions(param, fazzspec.CO_AND)
	return q
}

func (q *Query) OrGroupWhere(conditionFunc func(query *Query) *Query) *Query {
	query := QueryTx(q.Tx).Use(q.Model)
	param := conditionFunc(query).Parameter
	q.appendGroupConditions(param, fazzspec.CO_OR)
	return q
}

func (q *Query) GroupBy(column string) *Query {
	q.appendGroupBy(column)
	return q
}

func (q *Query) OrderBy(key string, direction fazzorder.OrderDirection) *Query {
	q.appendOrderBy(q.Model.GetTable(), key, direction)
	return q
}

func (q *Query) WithLimit(limit int) *Query {
	q.setLimit(limit)
	return q
}

func (q *Query) WithOffset(offset int) *Query {
	q.setOffset(offset)
	return q
}

func (q *Query) WithLock(lock fazzspec.Lock) *Query {
	q.setLock(lock)
	return q
}

func (q *Query) AppendCondition(connector fazzspec.Connector, key string, operator fazzspec.Operator, value interface{}) *Query {
	q.appendCondition(q.Model.GetTable(), connector, key, operator, value)
	return q
}