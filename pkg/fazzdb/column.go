package fazzdb

import "fmt"

// Col is a constructor for creating plain column
func Col(column string) Column {
	return Column{
		Key:       column,
		Aggregate: AG_NONE,
	}
}

// Avg is a constructor for creating avg column
func Avg(column string) Column {
	return Column{
		Key:       column,
		Aggregate: AG_AVG,
	}
}

// Count is a constructor for creating count column
func Count(column string) Column {
	return Column{
		Key:       column,
		Aggregate: AG_COUNT,
	}
}

// Sum is a constructor for creating sum column
func Sum(column string) Column {
	return Column{
		Key:       column,
		Aggregate: AG_SUM,
	}
}

// Min is a constructor for creating min column
func Min(column string) Column {
	return Column{
		Key:       column,
		Aggregate: AG_MIN,
	}
}

// Max is a constructor for creating max column
func Max(column string) Column {
	return Column{
		Key:       column,
		Aggregate: AG_MAX,
	}
}

// Column is a struct that is used to contain key and aggregate attribute for each column
type Column struct {
	Key       string
	Aggregate Aggregate
}

// ToString is a function that wrap toString function
func (c *Column) ToString(table string) string {
	return toString(table, *c)
}

// Order is a struct that is used to contain order by attributes
type Order struct {
	Field     Column
	Direction OrderDirection
	NullsLast bool
}

// ToString is a function that wrap toString function
func (o *Order) ToString(table string) string {
	return toString(table, o.Field)
}

// toString is a function that will handle building column query by given parameter
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
