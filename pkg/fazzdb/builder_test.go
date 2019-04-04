package fazzdb

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBuilder_BuildInsert(t *testing.T) {
	t.Run("BuildInsert AutoIncrementModel", func(t *testing.T) {
		aiQuery := builder.BuildInsert(ai, false)
		require.Equal(t, `INSERT INTO auto_increment_tests ( "name", "number" ) VALUES ( :name, :number ) RETURNING id;`, aiQuery)
	})

	t.Run("BuildInsert UuidModel", func(t *testing.T) {
		uQuery := builder.BuildInsert(u, false)
		require.Equal(t, `INSERT INTO uuid_tests ( "id", "name", "number" ) VALUES ( :id, :name, :number ) RETURNING id;`, uQuery)
	})

	t.Run("BuildInsert PlainModel", func(t *testing.T) {
		pQuery := builder.BuildInsert(p, false)
		require.Equal(t, `INSERT INTO plain_tests ( "id", "name", "number" ) VALUES ( :id, :name, :number ) RETURNING id;`, pQuery)
	})

	t.Run("BuildInsert TimestampModel", func(t *testing.T) {
		tsQuery := builder.BuildInsert(ts, false)
		require.Equal(t, `INSERT INTO timestamp_tests ( "createdAt", "updatedAt" ) VALUES ( :createdAt, :updatedAt ) RETURNING id;`, tsQuery)
	})

	t.Run("BuildInsert SoftDeleteModel", func(t *testing.T) {
		sdQuery := builder.BuildInsert(sd, false)
		require.Equal(t, `INSERT INTO soft_delete_tests ( "deletedAt" ) VALUES ( :deletedAt ) RETURNING id;`, sdQuery)
	})

	t.Run("BuildInsert CompleteModel", func(t *testing.T) {
		cQuery := builder.BuildInsert(c, false)
		require.Equal(t, `INSERT INTO complete_tests ( "name", "number", "createdAt", "updatedAt", "deletedAt" ) VALUES ( :name, :number, :createdAt, :updatedAt, :deletedAt ) RETURNING id;`, cQuery)
	})
}

func TestBuilder_BuildUpdate(t *testing.T) {
	t.Run("BuildUpdate AutoIncrementModel", func(t *testing.T) {
		param.Conditions = generateWhereIdParameter(ai)
		aiQuery := builder.BuildUpdate(ai, param)
		require.Equal(t, `UPDATE auto_increment_tests SET "name" = :name, "number" = :number WHERE  "auto_increment_tests"."id" = :id;`, aiQuery)
	})

	t.Run("BuildUpdate UuidModel", func(t *testing.T) {
		param.Conditions = generateWhereIdParameter(u)
		uQuery := builder.BuildUpdate(u, param)
		require.Equal(t, `UPDATE uuid_tests SET "name" = :name, "number" = :number WHERE  "uuid_tests"."id" = :id;`, uQuery)
	})

	t.Run("BuildUpdate PlainModel", func(t *testing.T) {
		param.Conditions = generateWhereIdParameter(p)
		pQuery := builder.BuildUpdate(p, param)
		require.Equal(t, `UPDATE plain_tests SET "name" = :name, "number" = :number WHERE  "plain_tests"."id" = :id;`, pQuery)
	})

	t.Run("BuildUpdate TimestampModel", func(t *testing.T) {
		param.Conditions = generateWhereIdParameter(ts)
		tsQuery := builder.BuildUpdate(ts, param)
		require.Equal(t, `UPDATE timestamp_tests SET "updatedAt" = :updatedAt WHERE  "timestamp_tests"."id" = :id;`, tsQuery)
	})

	t.Run("BuildUpdate SoftDeleteModel", func(t *testing.T) {
		param.Conditions = generateWhereIdParameter(sd)
		sdQuery := builder.BuildUpdate(sd, param)
		require.Equal(t, `UPDATE soft_delete_tests SET "deletedAt" = :deletedAt WHERE  "soft_delete_tests"."id" = :id;`, sdQuery)
	})

	t.Run("BuildUpdate CompleteModel", func(t *testing.T) {
		param.Conditions = generateWhereIdParameter(c)
		cQuery := builder.BuildUpdate(c, param)
		require.Equal(t, `UPDATE complete_tests SET "name" = :name, "number" = :number, "updatedAt" = :updatedAt, "deletedAt" = :deletedAt WHERE  "complete_tests"."id" = :id;`, cQuery)
	})
}

func TestBuilder_BuildDelete(t *testing.T) {
	t.Run("BuildDelete AutoIncrementModel", func(t *testing.T) {
		param.Conditions = generateWhereIdParameter(ai)
		aiQuery := builder.BuildDelete(ai, param)
		require.Equal(t, `DELETE FROM auto_increment_tests WHERE  "auto_increment_tests"."id" = :id;`, aiQuery)
	})

	t.Run("BuildDelete UuidModel", func(t *testing.T) {
		param.Conditions = generateWhereIdParameter(u)
		uQuery := builder.BuildDelete(u, param)
		require.Equal(t, `DELETE FROM uuid_tests WHERE  "uuid_tests"."id" = :id;`, uQuery)
	})

	t.Run("BuildDelete PlainModel", func(t *testing.T) {
		param.Conditions = generateWhereIdParameter(p)
		pQuery := builder.BuildDelete(p, param)
		require.Equal(t, `DELETE FROM plain_tests WHERE  "plain_tests"."id" = :id;`, pQuery)
	})

	t.Run("BuildDelete TimestampModel", func(t *testing.T) {
		param.Conditions = generateWhereIdParameter(ts)
		tsQuery := builder.BuildDelete(ts, param)
		require.Equal(t, `DELETE FROM timestamp_tests WHERE  "timestamp_tests"."id" = :id;`, tsQuery)
	})

	t.Run("BuildDelete SoftDeleteModel", func(t *testing.T) {
		param.Conditions = generateWhereIdParameter(sd)
		sdQuery := builder.BuildDelete(sd, param)
		require.Equal(t, `DELETE FROM soft_delete_tests WHERE  "soft_delete_tests"."id" = :id;`, sdQuery)
	})

	t.Run("BuildDelete CompleteModel", func(t *testing.T) {
		param.Conditions = generateWhereIdParameter(c)
		cQuery := builder.BuildDelete(c, param)
		require.Equal(t, `DELETE FROM complete_tests WHERE  "complete_tests"."id" = :id;`, cQuery)
	})
}

func TestBuilder_BuildBulkInsert(t *testing.T) {
	data := []interface{}{
		&CpTest{
			Model:  c.Model,
			Name:   "test_1",
			Number: 10,
		},
		&CpTest{
			Model:  c.Model,
			Name:   "test_2",
			Number: 11,
		},
	}

	cQuery := builder.BuildBulkInsert(c, data)
	require.Equal(t, `INSERT INTO complete_tests ( "name", "number", "createdAt", "updatedAt", "deletedAt" ) VALUES ( :0name, :0number, :0createdAt, :0updatedAt, :0deletedAt ), ( :1name, :1number, :1createdAt, :1updatedAt, :1deletedAt );`, cQuery)
}

func TestBuilder_BuildSelect(t *testing.T) {
	param.Conditions = []Condition{
		{
			Field:     Col("name"),
			Prefix:    "name0",
			Operator:  OP_LIKE,
			Connector: CO_NONE,
		},
		{
			Field:     Col("number"),
			Prefix:    "number0",
			Operator:  OP_LESS_THAN,
			Connector: CO_AND,
		},
		{
			Field:      Col("number"),
			Prefix:    "number1",
			Operator:  OP_MORE_THAN,
			Connector: CO_AND,
		},
		{
			Connector: CO_OR,
			Conditions: []Condition{
				{
					Field:     Col("id"),
					Prefix:    "id",
					Operator:  OP_EQUALS,
					Connector: CO_NONE,
				},
				{
					Field:     Col("createdAt"),
					Prefix:    "createdAt",
					Operator:  OP_LESS_THAN_EQUALS,
					Connector: CO_AND,
				},
				{
					Field:     Col("name"),
					Prefix:    "name1",
					Operator:  OP_EQUALS,
					Connector: CO_OR,
				},
			},
		},
		{
			Field:     Col("name"),
			Prefix:    "name2",
			Operator:  OP_IN,
			Connector: CO_AND,
		},
	}

	cQuery := builder.BuildSelect(c, param, AG_NONE, "")
	require.Equal(t, `SELECT  "complete_tests"."id", "complete_tests"."name", "complete_tests"."number", "complete_tests"."createdAt", "complete_tests"."updatedAt", "complete_tests"."deletedAt" FROM complete_tests WHERE  "complete_tests"."name" LIKE :name0 AND "complete_tests"."number" < :number0 AND "complete_tests"."number" > :number1 OR (  "complete_tests"."id" = :id AND "complete_tests"."createdAt" <= :createdAt OR "complete_tests"."name" = :name1 ) AND "complete_tests"."name" IN (:name2) LIMIT 15 OFFSET 5 FOR SHARE;`, cQuery)
}

func generateWhereIdParameter(m ModelInterface) []Condition {
	return []Condition{
		{
			Field:     Col(m.GetPK()),
			Prefix:    m.GetPK(),
			Operator:  OP_EQUALS,
			Connector: CO_NONE,
		},
	}
}
