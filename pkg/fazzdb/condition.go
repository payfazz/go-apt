package fazzdb

import (
	"fmt"
)

// SliceCondition is a struct that will handle condition in a slice
type SliceCondition struct {
	Connector  Connector
	Field      string
	Operator   Operator
	Value      interface{}
	Conditions []SliceCondition
}

// Condition is a struct that will handle condition when building query
type Condition struct {
	Field      Column
	Prefix     string
	Operator   Operator
	Connector  Connector
	Conditions []Condition
}

// QueryString is a function to build query based on given attributes
func (c *Condition) QueryString(table string) string {
	if len(c.Conditions) > 0 {
		var query = fmt.Sprintf("%s (", c.Connector)
		for _, cond := range c.Conditions {
			query = fmt.Sprintf("%s %s", query, cond.namedString(table))
		}
		query = fmt.Sprintf("%s )", query)
		return query
	}

	return c.namedString(table)
}

// namedString is a function to build condition query based on different operator
func (c *Condition) namedString(table string) string {
	query := ""
	key := c.Field.ToString(table)
	switch c.Operator {
	case OP_IS_NOT_NULL, OP_IS_NULL:
		query = fmt.Sprintf("%s %s %s", c.Connector, key, c.Operator)
	case OP_IN:
		query = fmt.Sprintf("%s %s %s (:%s)", c.Connector, key, c.Operator, c.Prefix)
	default:
		query = fmt.Sprintf("%s %s %s :%s", c.Connector, key, c.Operator, c.Prefix)
	}
	return query
}
