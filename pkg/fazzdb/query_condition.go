package fazzdb

// Where is a function that will add new condition that will check if column equals value given, connector
// that is used between condition is AND connector
func (q *Query) Where(key string, value interface{}) *Query {
	return q.AppendCondition(CO_AND, key, OP_EQUALS, value)
}

// WhereOp is a function that will add new condition that will check if column fulfill operator with value given,
// connector that is used between condition is AND connector
func (q *Query) WhereOp(key string, operator Operator, value interface{}) *Query {
	return q.AppendCondition(CO_AND, key, operator, value)
}

// WhereNil is a function that will add new condition that will check if column nil, connector that is used
// between condition is AND connector
func (q *Query) WhereNil(key string) *Query {
	return q.AppendCondition(CO_AND, key, OP_IS_NULL, nil)
}

// WhereNotNil is a function that will add new condition that will check if column not nil, connector that is used
// between condition is AND connector
func (q *Query) WhereNotNil(key string) *Query {
	return q.AppendCondition(CO_AND, key, OP_IS_NOT_NULL, nil)
}

// OrWhere is a function that will add new condition that will check if column equals value given, connector
// that is used between condition is OR connector
func (q *Query) OrWhere(key string, value interface{}) *Query {
	return q.AppendCondition(CO_OR, key, OP_EQUALS, value)
}

// OrWhereOp is a function that will add new condition that will check if column fulfill operator with value given,
// connector that is used between condition is OR connector
func (q *Query) OrWhereOp(key string, operator Operator, value interface{}) *Query {
	return q.AppendCondition(CO_OR, key, operator, value)
}

// OrWhereNil is a function that will add new condition that will check if column nil, connector that is used
// between condition is OR connector
func (q *Query) OrWhereNil(key string) *Query {
	return q.AppendCondition(CO_OR, key, OP_IS_NULL, nil)
}

// OrWhereNotNil is a function that will add new condition that will check if column not nil, connector that is used
// between condition is OR connector
func (q *Query) OrWhereNotNil(key string) *Query {
	return q.AppendCondition(CO_OR, key, OP_IS_NOT_NULL, nil)
}

// GroupWhere is a function that will receive a function that return a group of condition to be grouped
// together, connector that is used between condition is AND connector
func (q *Query) GroupWhere(conditionFunc func(query *Query) *Query) *Query {
	query := QueryTx(q.Tx, q.Config).Use(q.Model)
	param := conditionFunc(query).Parameter
	q.appendGroupConditions(param, CO_AND)
	return q
}

// OrGroupWhere is a function that will receive a function that return a group of condition to be grouped
// together, connector that is used between condition is OR connector
func (q *Query) OrGroupWhere(conditionFunc func(query *Query) *Query) *Query {
	query := QueryTx(q.Tx, q.Config).Use(q.Model)
	param := conditionFunc(query).Parameter
	q.appendGroupConditions(param, CO_OR)
	return q
}

// GroupBy is a function that will add new group by column
func (q *Query) GroupBy(column string) *Query {
	q.appendGroupBy(column)
	return q
}

// OrderBy is a function that will add new order by column with direction
func (q *Query) OrderBy(key string, direction OrderDirection) *Query {
	q.appendOrderBy(q.Model.GetTable(), key, direction)
	return q
}

// WithLimit is a function that will add limit to result row
func (q *Query) WithLimit(limit int) *Query {
	q.setLimit(limit)
	return q
}

// WithOffset is a function that will add offset to result row
func (q *Query) WithOffset(offset int) *Query {
	q.setOffset(offset)
	return q
}

// WithLock is a function that will add lock to current query
func (q *Query) WithLock(lock Lock) *Query {
	q.setLock(lock)
	return q
}

// AppendCondition is a wrapper function for appendCondition
func (q *Query) AppendCondition(connector Connector, key string, operator Operator, value interface{}) *Query {
	q.appendCondition(q.Model.GetTable(), connector, key, operator, value)
	return q
}