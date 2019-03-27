package fazzdb

import (
	"fmt"
)

type Builder struct {}

// BuildDelete is a function that will return delete query from given model and parameter
func (b *Builder) BuildDelete(m ModelInterface, param *Parameter) string {
	query := fmt.Sprintf("DELETE FROM %s", m.GetTable())
	query = b.generateConditions(query, param)
	query = fmt.Sprintf("%s;", query)
	return query
}

// BuildUpdate is a function that will return delete query from given model and parameter
func (b *Builder) BuildUpdate(m ModelInterface, param *Parameter) string {
	query := fmt.Sprintf("UPDATE %s SET", m.GetTable())
	query = b.generateValues(query, m, b.isPrimaryKey, b.generateUpdateColumns)
	query = b.generateConditions(query, param)
	query = fmt.Sprintf("%s;", query)
	return query
}

// BuildBulkInsert is a function that will return bulk insert query from given model and slice of data
func (b *Builder) BuildBulkInsert(m ModelInterface, data []interface{}) string {
	query := fmt.Sprintf("INSERT INTO %s", m.GetTable())

	query = fmt.Sprintf("%s (", query)
	query = b.generateValues(query, m, b.isAutoIncrementPrimaryKey, b.generateSelectColumns)
	query = fmt.Sprintf("%s ) VALUES", query)

	firstData := true
	for i, v := range data {
		if firstData {
			query = fmt.Sprintf("%s (", query)
			firstData = false
		} else {
			query = fmt.Sprintf("%s, (", query)
		}

		query = b.generateBulkValues(query, v.(ModelInterface), i)

		query = fmt.Sprintf("%s )", query)
	}

	return query
}

// BuildInsert is a function that will return insert query from given model
func (b *Builder) BuildInsert(m ModelInterface) string {
	query := fmt.Sprintf("INSERT INTO %s", m.GetTable())

	query = fmt.Sprintf("%s (", query)
	query = b.generateValues(query, m, b.isAutoIncrementPrimaryKey, b.generateSelectColumns)

	query = fmt.Sprintf("%s ) VALUES (", query)
	query = b.generateValues(query, m, b.isAutoIncrementPrimaryKey, b.generateInsertValues)

	query = fmt.Sprintf("%s ) RETURNING %s;", query, m.GetPK())
	return query
}

// BuildSelect is a function that will return select query from given model and parameter
func (b *Builder) BuildSelect(m ModelInterface, param *Parameter, aggregate Aggregate, aggregateColumn string) string {
	query := "SELECT "

	if aggregate != AG_NONE {
		query = fmt.Sprintf("%s %s(%s)", query, aggregate, aggregateColumn)
	} else if m.ColumnCount() != 0 {
		query = b.generateValues(query, m, b.alwaysFalse, b.generateSelectColumns)
	} else {
		query = fmt.Sprintf("%s *", query)
	}

	query = fmt.Sprintf("%s FROM %s", query, m.GetTable())
	query = b.generateConditions(query, param)

	if len(param.Orders) > 0 {
		query = fmt.Sprintf("%s ORDER BY", query)
		for i, order := range param.Orders {
			if i == 0 {
				query = fmt.Sprintf("%s %s.%s %s", query, order.Table, order.Key, order.Direction)
			} else {
				query = fmt.Sprintf("%s, %s.%s %s", query, order.Table, order.Key, order.Direction)
			}
		}
	}

	if len(param.Groups) > 0 {
		query = fmt.Sprintf("%s GROUP BY", query)
		for i, group := range param.Groups {
			if i == 0 {
				query = fmt.Sprintf("%s %s.%s", query, m.GetTable(), group)
			} else {
				query = fmt.Sprintf("%s, %s.%s", query, m.GetTable(), group)
			}
		}
	}

	if param.Limit > 0 {
		query = fmt.Sprintf("%s LIMIT %d", query, param.Limit)
	}

	if param.Offset > 0 {
		query = fmt.Sprintf("%s OFFSET %d", query, param.Offset)
	}

	query = fmt.Sprintf("%s %s;", query, param.Lock)
	return query
}

// alwaysFalse is a function that will always skip id column when building query
func (b *Builder) alwaysFalse(column string, m ModelInterface) bool {
	return false
}

// isAutoIncrementPrimaryKey is a function that will skip id column if model is autoincrement when building query
func (b *Builder) isAutoIncrementPrimaryKey(column string, m ModelInterface) bool {
	return column == m.GetPK() && !m.IsUuid() && m.IsAutoIncrement()
}

// isPrimaryKey is a function that will skip id column when building query
func (b *Builder) isPrimaryKey(column string, m ModelInterface) bool {
	return column == m.GetPK()
}

// generateInsertValues is a function that will generate insert arguments for query
func (b *Builder) generateInsertValues(query string, column string, first bool) (string, bool) {
	if first {
		query = fmt.Sprintf("%s :%s", query, column)
		first = false
	} else {
		query = fmt.Sprintf("%s, :%s", query, column)
	}
	return query, first
}

// generateUpdateColumns is a function that will generate update column with arguments for query
func (b *Builder) generateUpdateColumns(query string, column string, first bool) (string, bool) {
	if first {
		query = fmt.Sprintf("%s \"%s\" = :%s", query, column, column)
		first = false
	} else {
		query = fmt.Sprintf("%s, \"%s\" = :%s", query, column, column)
	}
	return query, first
}

// generateSelectColumns is a function that will generate select columns
func (b *Builder) generateSelectColumns(query string, column string, first bool) (string, bool) {
	if first {
		query = fmt.Sprintf("%s \"%s\"", query, column)
		first = false
	} else {
		query = fmt.Sprintf("%s, \"%s\"", query, column)
	}
	return query, first
}

// generateBulkValues is a function that will generate insert arguments for query sequentially
func (b *Builder) generateBulkValues(query string, model ModelInterface, index int) string {
	first := true
	for _, column := range model.GetColumns() {
		if b.isAutoIncrementPrimaryKey(column, model) {
			continue
		}
		query, first = b.generateInsertValues(query, fmt.Sprintf("%d%s", index, column), first)
	}
	return query
}

// generateValues is a wrapper function for calling generate*Columns or generate*Values
func (b *Builder) generateValues(
	query string,
	m ModelInterface,
	skipFunc func(column string, m ModelInterface) bool,
	generate func(query string, column string, first bool) (string, bool),
) string {
	first := true
	for _, column := range m.GetColumns() {
		if skipFunc(column, m) {
			continue
		}
		query, first = generate(query, column, first)
	}
	return query
}

// generateConditions is a function that will generate condition based on given parameter
func (b *Builder) generateConditions(query string, param *Parameter) string {
	if len(param.Conditions) > 0 {
		param.Conditions[0].Connector = CO_EMPTY

		query = fmt.Sprintf("%s WHERE", query)
		for _, cond := range param.Conditions {
			query = fmt.Sprintf("%s %s", query, cond.QueryString())
		}
	}
	return query
}