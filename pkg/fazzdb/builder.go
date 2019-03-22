package fazzdb

import (
	"fmt"
)

type Builder struct {}

func (b *Builder) BuildDelete(m ModelInterface, param *Parameter) string {
	// SET TABLE
	query := fmt.Sprintf("DELETE FROM %s", m.GetTable())

	// SET CONDITION
	query = b.generateConditions(query, param)

	// END QUERY
	query = fmt.Sprintf("%s;", query)
	return query
}

func (b *Builder) BuildUpdate(m ModelInterface, param *Parameter) string {
	// SET TABLE
	query := fmt.Sprintf("UPDATE %s SET", m.GetTable())

	// SET COLUMN
	query = b.generateValues(query, m, b.isPrimaryKey, b.generateUpdateColumns)

	// SET CONDITION
	query = b.generateConditions(query, param)

	// END QUERY
	query = fmt.Sprintf("%s;", query)
	return query
}

func (b *Builder) BuildInsert(m ModelInterface) string {
	// SET TABLE
	query := fmt.Sprintf("INSERT INTO %s", m.GetTable())

	// SET COLUMN
	query = fmt.Sprintf("%s (", query)
	query = b.generateValues(query, m, b.isAutoIncrementPrimaryKey, b.generateSelectColumns)
	query = fmt.Sprintf("%s ) VALUES (", query)

	// SET NAMED VALUE
	query = b.generateValues(query, m, b.isAutoIncrementPrimaryKey, b.generateInsertValues)

	// SET RETURNING
	query = fmt.Sprintf("%s ) RETURNING %s;", query, m.GetPK())
	return query
}

func (b *Builder) BuildSelect(m ModelInterface, param *Parameter) string {
	query := "SELECT "

	// SET TABLE COLUMNS
	if m.ColumnCount() != 0 {
		query = b.generateValues(query, m, b.alwaysFalse, b.generateSelectColumns)
	} else {
		query = fmt.Sprintf("%s *", query)
	}

	// SET TABLE REFERENCE
	query = fmt.Sprintf("%s FROM %s", query, m.GetTable())

	// SET CONDITIONS
	query = b.generateConditions(query, param)

	// SET ORDER BY
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

	// SET GROUP BY
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

	// SET LIMIT
	if param.Limit > 0 {
		query = fmt.Sprintf("%s LIMIT %d", query, param.Limit)
	}

	// SET OFFSET
	if param.Offset > 0 {
		query = fmt.Sprintf("%s OFFSET %d", query, param.Offset)
	}

	// SET LOCK
	query = fmt.Sprintf("%s %s;", query, param.Lock)
	return query
}

func (b *Builder) alwaysFalse(column string, m ModelInterface) bool {
	return false
}

func (b *Builder) isAutoIncrementPrimaryKey(column string, m ModelInterface) bool {
	return column == m.GetPK() && !m.IsUuid() && m.IsAutoIncrement()
}

func (b *Builder) isPrimaryKey(column string, m ModelInterface) bool {
	return column == m.GetPK()
}

func (b *Builder) generateInsertValues(query string, column string, first bool) (string, bool) {
	if first {
		query = fmt.Sprintf("%s :%s", query, column)
		first = false
	} else {
		query = fmt.Sprintf("%s, :%s", query, column)
	}
	return query, first
}

func (b *Builder) generateUpdateColumns(query string, column string, first bool) (string, bool) {
	if first {
		query = fmt.Sprintf("%s %s = :%s", query, column, column)
		first = false
	} else {
		query = fmt.Sprintf("%s, %s = :%s", query, column, column)
	}
	return query, first
}

func (b *Builder) generateSelectColumns(query string, column string, first bool) (string, bool) {
	if first {
		query = fmt.Sprintf("%s %s", query, column)
		first = false
	} else {
		query = fmt.Sprintf("%s, %s", query, column)
	}
	return query, first
}

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