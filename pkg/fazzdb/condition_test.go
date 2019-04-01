package fazzdb

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCondition_QueryString(t *testing.T) {
	single := "single_tests"
	group := "group_tests"

	// OR SINGLE EQUALS
	cond := Condition{
		Field: Col("name"),
		Prefix: "name",
		Operator: OP_EQUALS,
		Connector: CO_OR,
	}

	result := cond.QueryString(single)
	require.Equal(t, "OR \"single_tests\".\"name\" = :name", result)

	// AND SINGLE IS NULL
	cond = Condition{
		Field: Col("name"),
		Prefix: "name",
		Operator: OP_IS_NULL,
		Connector: CO_AND,
	}

	result = cond.QueryString(single)
	require.Equal(t, "AND \"single_tests\".\"name\" IS NULL", result)

	// NONE SINGLE IN
	cond = Condition{
		Field: Col("name"),
		Prefix: "name",
		Operator: OP_IN,
		Connector: CO_NONE,
	}

	result = cond.QueryString(single)
	require.Equal(t, " \"single_tests\".\"name\" IN (:name)", result)

	// AND GROUP
	cond = Condition{
		Connector: CO_AND,
		Conditions: []Condition{
			{
				Field: Col("name"),
				Prefix: "name0",
				Operator: OP_EQUALS,
				Connector: CO_NONE,
			},
			{
				Field: Col("name"),
				Prefix: "name1",
				Operator: OP_LESS_THAN,
				Connector: CO_AND,
			},
			{
				Field: Col("name"),
				Prefix: "name2",
				Operator: OP_NOT_LIKE,
				Connector: CO_OR,
			},
		},
	}

	result = cond.QueryString(group)
	require.Equal(t, "AND (  \"group_tests\".\"name\" = :name0 AND \"group_tests\".\"name\" < :name1 OR \"group_tests\".\"name\" NOT LIKE :name2 )", result)
}
