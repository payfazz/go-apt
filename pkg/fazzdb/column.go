package fazzdb

import "fmt"

func Col(column string) Column {
	return Column{
		Key:       column,
		Aggregate: AG_NONE,
	}
}

func Avg(column string) Column {
	return Column{
		Key:       column,
		Aggregate: AG_AVG,
	}
}

func Count(column string) Column {
	return Column{
		Key:       column,
		Aggregate: AG_COUNT,
	}
}

func Sum(column string) Column {
	return Column{
		Key:       column,
		Aggregate: AG_SUM,
	}
}

func Min(column string) Column {
	return Column{
		Key:       column,
		Aggregate: AG_MIN,
	}
}

func Max(column string) Column {
	return Column{
		Key:       column,
		Aggregate: AG_MAX,
	}
}

type Column struct {
	Key       string
	Aggregate Aggregate
}

func (c *Column) ToString(table string) string {
	return toString(table, *c)
}

// Order is a struct that is used to contain order by attributes
type Order struct {
	Field     Column
	Direction OrderDirection
}

func (o *Order) ToString(table string) string {
	return toString(table, o.Field)
}

func toString(table string, c Column) string {
	result := fmt.Sprintf(`"%s"."%s"`, table, c.Key)
	if "" == table {
		result = fmt.Sprintf(`"%s"`, c.Key)
	}
	if AG_NONE != c.Aggregate {
		result = fmt.Sprintf(`%s("%s")`, c.Aggregate, c.Key)
	}
	return result
}
