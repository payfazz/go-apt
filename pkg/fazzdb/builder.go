package fazzdb

import (
	`fmt`
	`sync`
)

var once sync.Once
var singleton *Builder

// NewBuilder is a constructor to retrieve instance of builder
func NewBuilder() *Builder {
	once.Do(func() {
		singleton = &Builder{}
	})

	return singleton
}

// Builder is a struct that will handle transforming parameters into query string
type Builder struct {}

// BuildDelete is a function that will return delete query from given model and parameter
func (b *Builder) BuildDelete(model ModelInterface, param *Parameter) string {
	query := fmt.Sprintf(`DELETE FROM %s`, model.GetTable())
	query = b.generateConditions(query, model.GetTable(), param)
	query = fmt.Sprintf(`%s;`, query)
	return query
}

// BuildUpdate is a function that will return delete query from given model and parameter
func (b *Builder) BuildUpdate(model ModelInterface, param *Parameter) string {
	query := fmt.Sprintf(`UPDATE %s SET`, model.GetTable())
	query = b.generateValues(query, model, param, b.isPrimaryKeyOrCreatedAt, b.generateUpdateColumns)
	query = b.generateConditions(query, model.GetTable(), param)
	query = fmt.Sprintf(`%s;`, query)
	return query
}

// BuildBulkInsert is a function that will return bulk insert query from given model and slice of data
func (b *Builder) BuildBulkInsert(model ModelInterface, data []interface{}) string {
	query := fmt.Sprintf(`INSERT INTO %s`, model.GetTable())

	query = fmt.Sprintf(`%s (`, query)
	query = b.generateValues(query, model, nil, b.isAutoIncrementPrimaryKey, b.generateInsertColumns)
	query = fmt.Sprintf(`%s ) VALUES`, query)

	firstData := true
	for i, v := range data {
		if firstData {
			query = fmt.Sprintf(`%s (`, query)
			firstData = false
		} else {
			query = fmt.Sprintf(`%s, (`, query)
		}

		query = b.generateBulkValues(query, v.(ModelInterface), i)

		query = fmt.Sprintf(`%s )`, query)
	}

	query = fmt.Sprintf(`%s;`, query)
	return query
}

// BuildInsert is a function that will return insert query from given model
func (b *Builder) BuildInsert(model ModelInterface) string {
	query := fmt.Sprintf(`INSERT INTO %s`, model.GetTable())

	query = fmt.Sprintf(`%s (`, query)
	query = b.generateValues(query, model, nil, b.isAutoIncrementPrimaryKey, b.generateInsertColumns)

	query = fmt.Sprintf(`%s ) VALUES (`, query)
	query = b.generateValues(query, model, nil, b.isAutoIncrementPrimaryKey, b.generateInsertValues)

	query = fmt.Sprintf(`%s ) RETURNING %s;`, query, model.GetPK())
	return query
}

// BuildSelect is a function that will return select query from given model and parameter
func (b *Builder) BuildSelect(model ModelInterface, param *Parameter, aggregate Aggregate, aggregateColumn string) string {
	query := `SELECT `

	if aggregate != AG_NONE {
		query = fmt.Sprintf(`%s %s(%s)`, query, aggregate, aggregateColumn)
	} else if model.ColumnCount() != 0 {
		query = b.generateValues(query, model, param, b.alwaysFalse, b.generateSelectColumns)
	} else {
		query = fmt.Sprintf(`%s *`, query)
	}

	query = fmt.Sprintf(`%s FROM %s`, query, model.GetTable())
	query = b.generateConditions(query, model.GetTable(), param)

	if len(param.Groups) > 0 {
		query = fmt.Sprintf(`%s GROUP BY`, query)
		for i, group := range param.Groups {
			key := group.ToString(model.GetTable())
			if i == 0 {
				query = fmt.Sprintf(`%s %s`, query, key)
			} else {
				query = fmt.Sprintf(`%s, %s`, query, key)
			}
		}
	}

	query = b.generateHavingConditions(query, model.GetTable(), param)

	if len(param.Orders) > 0 {
		query = fmt.Sprintf(`%s ORDER BY`, query)
		for i, order := range param.Orders {
			key := order.Field.ToString(model.GetTable())
			if i == 0 {
				query = fmt.Sprintf(`%s %s %s`, query, key, order.Direction)
			} else {
				query = fmt.Sprintf(`%s, %s %s`, query, key, order.Direction)
			}
		}
	}

	if param.Limit > 0 {
		query = fmt.Sprintf(`%s LIMIT %d`, query, param.Limit)
	}

	if param.Offset > 0 {
		query = fmt.Sprintf(`%s OFFSET %d`, query, param.Offset)
	}

	query = fmt.Sprintf(`%s %s;`, query, param.Lock)
	return query
}

// alwaysFalse is a function that will always skip id column when building query
func (b *Builder) alwaysFalse(column string, model ModelInterface) bool {
	return false
}

// isAutoIncrementPrimaryKey is a function that will skip id column if model is autoincrement when building query
func (b *Builder) isAutoIncrementPrimaryKey(column string, model ModelInterface) bool {
	return column == model.GetPK() && !model.IsUuid() && model.IsAutoIncrement()
}

// isPrimaryKey is a function that will skip id column when building query
func (b *Builder) isPrimaryKeyOrCreatedAt(column string, model ModelInterface) bool {
	return column == model.GetPK() || column == CREATED_AT
}

// generateInsertValues is a function that will generate insert arguments for query
func (b *Builder) generateInsertValues(query string, table string, column Column, first bool) (string, bool) {
	if first {
		query = fmt.Sprintf(`%s :%s`, query, column.Key)
		first = false
	} else {
		query = fmt.Sprintf(`%s, :%s`, query, column.Key)
	}
	return query, first
}

// generateUpdateColumns is a function that will generate update column with arguments for query
func (b *Builder) generateUpdateColumns(query string, table string, column Column, first bool) (string, bool) {
	if first {
		query = fmt.Sprintf(`%s "%s" = :%s`, query, column.Key, column.Key)
		first = false
	} else {
		query = fmt.Sprintf(`%s, "%s" = :%s`, query, column.Key, column.Key)
	}
	return query, first
}

// generateSelectColumns is a function that will generate insert columns
func (b *Builder) generateInsertColumns(query string, table string, column Column, first bool) (string, bool) {
	if first {
		query = fmt.Sprintf(`%s "%s"`, query, column.Key)
		first = false
	} else {
		query = fmt.Sprintf(`%s, "%s"`, query, column.Key)
	}
	return query, first
}

// generateSelectColumns is a function that will generate select columns
func (b *Builder) generateSelectColumns(query string, table string, column Column, first bool) (string, bool) {
	col := column.ToString(table)
	if first {
		query = fmt.Sprintf(`%s %s`, query, col)
		first = false
	} else {
		query = fmt.Sprintf(`%s, %s`, query, col)
	}
	return query, first
}

// generateBulkValues is a function that will generate insert arguments for query sequentially
func (b *Builder) generateBulkValues(query string, model ModelInterface, index int) string {
	first := true
	for _, column := range model.GetColumns() {
		if b.isAutoIncrementPrimaryKey(column.Key, model) {
			continue
		}

		tempCol := column
		tempCol.Key = fmt.Sprintf(`%d%s`, index, column.Key)

		query, first = b.generateInsertValues(query, model.GetTable(), tempCol, first)
	}
	return query
}

// generateValues is a wrapper function for calling generate*Columns or generate*Values
func (b *Builder) generateValues(
	query string,
	model ModelInterface,
	param *Parameter,
	skipFunc func(column string, model ModelInterface) bool,
	generate func(query string, table string, column Column, first bool) (string, bool),
) string {
	first := true

	columns := model.GetColumns()
	if nil != param && len(param.Columns) > 0 {
		columns = param.Columns
	}

	for _, column := range columns {
		if skipFunc(column.Key, model) {
			continue
		}
		query, first = generate(query, model.GetTable(), column, first)
	}
	return query
}

// generateConditions is a function that will generate condition based on given parameter
func (b *Builder) generateConditions(query string, table string, param *Parameter) string {
	if len(param.Conditions) > 0 {
		param.Conditions[0].Connector = CO_NONE

		query = fmt.Sprintf(`%s WHERE`, query)
		for _, cond := range param.Conditions {
			query = fmt.Sprintf(`%s %s`, query, cond.QueryString(table))
		}
	}
	return query
}

// generateHavingConditions is a function that will generate having condition based on given parameter
func (b *Builder) generateHavingConditions(query string, table string, param *Parameter) string {
	if len(param.Havings) > 0 {
		param.Havings[0].Connector = CO_NONE

		query = fmt.Sprintf(`%s HAVING`, query)
		for _, cond := range param.Havings {
			query = fmt.Sprintf(`%s %s`, query, cond.QueryString(table))
		}
	}
	return query
}