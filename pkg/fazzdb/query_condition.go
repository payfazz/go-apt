package fazzdb

func (q *Query) Where(key string, value interface{}) *Query {
	return q.AppendCondition(CO_AND, key, OP_EQUALS, value)
}

func (q *Query) WhereOp(key string, operator Operator, value interface{}) *Query {
	return q.AppendCondition(CO_AND, key, operator, value)
}

func (q *Query) OrWhere(key string, value interface{}) *Query {
	return q.AppendCondition(CO_OR, key, OP_EQUALS, value)
}

func (q *Query) OrWhereOp(key string, operator Operator, value interface{}) *Query {
	return q.AppendCondition(CO_OR, key, operator, value)
}

func (q *Query) GroupWhere(conditionFunc func(query *Query) *Query) *Query {
	query := QueryTx(q.Tx, q.Config).Use(q.Model)
	param := conditionFunc(query).Parameter
	q.appendGroupConditions(param, CO_AND)
	return q
}

func (q *Query) OrGroupWhere(conditionFunc func(query *Query) *Query) *Query {
	query := QueryTx(q.Tx, q.Config).Use(q.Model)
	param := conditionFunc(query).Parameter
	q.appendGroupConditions(param, CO_OR)
	return q
}

func (q *Query) GroupBy(column string) *Query {
	q.appendGroupBy(column)
	return q
}

func (q *Query) OrderBy(key string, direction OrderDirection) *Query {
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

func (q *Query) WithLock(lock Lock) *Query {
	q.setLock(lock)
	return q
}

func (q *Query) AppendCondition(connector Connector, key string, operator Operator, value interface{}) *Query {
	q.appendCondition(q.Model.GetTable(), connector, key, operator, value)
	return q
}