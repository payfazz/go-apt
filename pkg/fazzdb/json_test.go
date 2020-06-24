package fazzdb

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestJSONPathOp(t *testing.T) {
	tableName := "tests"

	// OR JSON KEY EQUALS
	cond := Condition{
		Field:     Col("data"),
		Prefix:    "name",
		Operator:  JSONFieldOp(JSONPath{"name"}, OP_EQUALS),
		Connector: CO_OR,
	}

	result := cond.QueryString(tableName)
	require.Equal(t, "OR \"tests\".\"data\" #>> '{name}' = :name", result)

	// AND JSON ARRAY IS NULL
	cond = Condition{
		Field:     Col("data"),
		Prefix:    "name",
		Operator:  JSONFieldOp(JSONPath{0}, OP_IS_NULL),
		Connector: CO_AND,
	}

	result = cond.QueryString(tableName)
	require.Equal(t, "AND \"tests\".\"data\" #>> '{0}' IS NULL", result)

	// NONE PATH KEY IN
	cond = Condition{
		Field:     Col("data"),
		Prefix:    "name",
		Operator:  JSONFieldOp(JSONPath{"child", "lists", 0}, OP_IN),
		Connector: CO_NONE,
	}

	result = cond.QueryString(tableName)
	require.Equal(t, " \"tests\".\"data\" #>> '{child,lists,0}' IN (:name)", result)
}
