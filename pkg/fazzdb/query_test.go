package fazzdb

import (
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func initQuery() *Query {
	conn := "host=localhost port=5432 user=postgres password=cashfazz dbname=fazzdb_test sslmode=disable"
	db, err := sqlx.Connect("postgres", conn)
	if nil != err {
		panic(err)
	}

	return QueryDb(db, Config{
		Limit:           0,
		Offset:          0,
		Lock:            LO_NONE,
		DevelopmentMode: true,
	})
}

func initTestDb(query *Query) {
	qs := `DROP SCHEMA public CASCADE; CREATE SCHEMA public;

create table complete_tests
(
  id        serial not null
    constraint complete_tests_pk
      primary key,
  name      varchar,
  number    integer,
  "created_at" timestamp,
  "updated_at" timestamp,
  "deleted_at" timestamp
);`

	_, err := query.RawExec(qs)
	if nil != err {
		panic(err)
	}
}

func initTestData(query *Query) {
	qd := `INSERT INTO complete_tests ("name", "number", "created_at", "updated_at", "deleted_at") VALUES
('test_name_1', 23, '2019-03-28 02:00:00', '2019-03-29 03:00:00', NULL),
('test_name_2', 5, '2019-03-28 02:00:00', NULL, NULL),
('test_name_3', 19, '2019-03-28 02:00:00', NULL, NULL),
('test_name_4', 31, '2019-03-28 02:00:00', NULL, NULL),
('test_name_5', 40, '2019-03-28 02:00:00', '2019-03-29 03:00:00', NULL),
('test_name_6', 36, '2019-03-28 02:00:00', NULL, NULL),
('test_name_7', 16, '2019-03-28 02:00:00', NULL, NULL),
('test_name_8', 21, '2019-03-28 02:00:00', NULL, '2019-03-30 02:00:00'),
('test_name_9', 26, '2019-03-28 02:00:00', NULL, NULL),
('test_name_10', 32, '2019-03-28 02:00:00', NULL, NULL),
('test_name_11', 13, '2019-03-28 02:00:00', NULL, NULL),
('test_name_12', 18, '2019-03-28 02:00:00', NULL, '2019-03-30 01:00:00');`

	_, err := query.RawExec(qd)
	if nil != err {
		panic(err)
	}
}

func initGroupByData(query *Query) {
	qd := `INSERT INTO complete_tests ("name", "number", "created_at") VALUES
('test_name_1', 23, '2019-03-28 02:00:00'),
('test_name_1', 5, '2019-03-28 02:00:00'),
('test_name_1', 19, '2019-03-28 02:00:00'),
('test_name_2', 31, '2019-03-28 02:00:00'),
('test_name_2', 40, '2019-03-28 02:00:00'),
('test_name_2', 36, '2019-03-28 02:00:00'),
('test_name_3', 16, '2019-03-28 02:00:00'),
('test_name_3', 21, '2019-03-28 02:00:00'),
('test_name_3', 26, '2019-03-28 02:00:00'),
('test_name_3', 32, '2019-03-28 02:00:00'),
('test_name_4', 13, '2019-03-28 02:00:00'),
('test_name_4', 18, '2019-03-28 02:00:00');`

	_, err := query.RawExec(qd)
	if nil != err {
		panic(err)
	}
}

func TestQuery_Exec(t *testing.T) {
	query := initQuery()
	initTestDb(query)

	t.Run("Insert", func(t *testing.T) {
		c.Name = "c_name_test"
		c.Number = 35

		id, err := query.Use(c).Insert()
		require.NoError(t, err)
		require.NotNil(t, id)

		row, err := query.Use(c).
			Where("id", id).
			First()
		require.NoError(t, err)
		require.NotNil(t, row)

		result := row.(*CpTest)

		require.Equal(t, 1, result.Id)
		require.Equal(t, "c_name_test", result.Name)
		require.Equal(t, 35, result.Number)
		require.NotNil(t, result.CreatedAt)
		require.Nil(t, result.UpdatedAt)
		require.Nil(t, result.DeletedAt)
	})

	t.Run("Update", func(t *testing.T) {
		row, err := query.Use(c).
			Where("id", 1).
			First()
		require.NoError(t, err)

		result := row.(*CpTest)
		result.Name = "c_update_name_test"
		result.Number = 25

		_, err = query.Use(result).Update()
		require.NoError(t, err)

		row, err = query.Use(c).
			Where("id", 1).
			First()
		require.NoError(t, err)
		require.NotNil(t, row)

		result = row.(*CpTest)

		require.Equal(t, 1, result.Id)
		require.Equal(t, "c_update_name_test", result.Name)
		require.Equal(t, 25, result.Number)
		require.NotNil(t, result.CreatedAt)
		require.NotNil(t, result.UpdatedAt)
		require.Nil(t, result.DeletedAt)
	})

	t.Run("Delete", func(t *testing.T) {
		row, err := query.Use(c).Where("id", 1).First()
		require.NoError(t, err)
		result := row.(*CpTest)

		success, err := query.Use(result).Delete()
		require.NoError(t, err)
		require.Equal(t, success, true)

		row, err = query.Use(c).Where("id", 1).First()
		require.Error(t, err)
		require.Nil(t, row)

		row, err = query.Use(c).Where("id", 1).FirstWithTrash()
		require.NoError(t, err)
		require.NotNil(t, row)

		require.Equal(t, 1, result.Id)
		require.Equal(t, "c_update_name_test", result.Name)
		require.Equal(t, 25, result.Number)
		require.NotNil(t, result.CreatedAt)
		require.NotNil(t, result.UpdatedAt)
		require.NotNil(t, result.DeletedAt)
	})
}

func TestQuery_Select(t *testing.T) {
	query := initQuery()
	initTestDb(query)
	initTestData(query)

	c = NewCpTest()

	t.Run("First", func(t *testing.T) {
		row, err := query.Use(c).First()
		require.NoError(t, err)
		require.NotNil(t, row)

		result := row.(*CpTest)
		require.Equal(t, 1, result.Id)
		require.Equal(t, "test_name_1", result.Name)
		require.Equal(t, 23, result.Number)
		require.NotNil(t, result.CreatedAt)
		require.NotNil(t, result.UpdatedAt)
		require.Nil(t, result.DeletedAt)

		row, err = query.Use(c).OrderBy("id", DIR_DESC).
			FirstWithTrash()
		require.NoError(t, err)
		require.NotNil(t, row)

		result = row.(*CpTest)
		require.Equal(t, 12, result.Id)
		require.Equal(t, "test_name_12", result.Name)
		require.Equal(t, 18, result.Number)
		require.NotNil(t, result.CreatedAt)
		require.Nil(t, result.UpdatedAt)
		require.NotNil(t, result.DeletedAt)
	})

	t.Run("All", func(t *testing.T) {
		row, err := query.Use(c).
			OrWhereNil(DELETED_AT).
			OrWhereNotNil(CREATED_AT).
			All()
		require.NoError(t, err)
		require.NotNil(t, row)

		result := row.([]*CpTest)

		require.Equal(t, 10, len(result))
		require.Equal(t, 11, result[9].Id)
		require.Equal(t, "test_name_11", result[9].Name)
		require.Equal(t, 13, result[9].Number)
		require.NotNil(t, result[9].CreatedAt)
		require.Nil(t, result[9].UpdatedAt)
		require.Nil(t, result[9].DeletedAt)

		row, err = query.Use(c).
			AllWithTrash()
		require.NoError(t, err)
		require.NotNil(t, row)

		result = row.([]*CpTest)

		require.Equal(t, 12, len(result))
		require.Equal(t, 10, result[9].Id)
		require.Equal(t, "test_name_10", result[9].Name)
		require.Equal(t, 32, result[9].Number)
		require.NotNil(t, result[9].CreatedAt)
		require.Nil(t, result[9].UpdatedAt)
		require.Nil(t, result[9].DeletedAt)
	})

	t.Run("Number Between", func(t *testing.T) {
		row, err := query.Use(c).WhereOp("number", OP_LESS_THAN, 25).
			WhereOp("number", OP_MORE_THAN, 10).
			All()
		require.NoError(t, err)
		require.NotNil(t, row)

		result := row.([]*CpTest)

		require.Equal(t, 4, len(result))
		require.Equal(t, 11, result[3].Id)
		require.Equal(t, "test_name_11", result[3].Name)
		require.Equal(t, 13, result[3].Number)
		require.NotNil(t, result[3].CreatedAt)
		require.Nil(t, result[3].UpdatedAt)
		require.Nil(t, result[3].DeletedAt)
	})

	t.Run("Where", func(t *testing.T) {
		row, err := query.Use(c).GroupWhere(func(query *Query) *Query {
			return query.Where("number", 23).
				OrWhereOp("number", OP_LESS_THAN, 5)
		}).
			All()
		require.NoError(t, err)
		require.NotNil(t, row)

		result := row.([]*CpTest)

		require.Equal(t, 1, len(result))
		require.Equal(t, 1, result[0].Id)
		require.Equal(t, "test_name_1", result[0].Name)
		require.Equal(t, 23, result[0].Number)
		require.NotNil(t, result[0].CreatedAt)
		require.NotNil(t, result[0].UpdatedAt)
		require.Nil(t, result[0].DeletedAt)

		inPayload := []interface{}{
			23, 5, 19, 32, 18,
		}

		cnt, err := query.Use(c).WhereOp("number", OP_IN, inPayload).
			Count()
		require.NoError(t, err)
		require.Equal(t, float64(4), *cnt)
	})

	t.Run("Order By", func(t *testing.T) {
		row, err := query.Use(c).OrderBy("number", DIR_DESC).
			All()
		require.NoError(t, err)
		require.NotNil(t, row)

		result := row.([]*CpTest)

		require.Equal(t, 10, len(result))
		require.Equal(t, 5, result[0].Id)
		require.Equal(t, "test_name_5", result[0].Name)
		require.Equal(t, 40, result[0].Number)
		require.NotNil(t, result[0].CreatedAt)
		require.NotNil(t, result[0].UpdatedAt)
		require.Nil(t, result[0].DeletedAt)
	})

	t.Run("Limit", func(t *testing.T) {
		row, err := query.Use(c).WithLimit(5).
			WithLock(LO_FOR_UPDATE).
			All()
		require.NoError(t, err)
		require.NotNil(t, row)

		result := row.([]*CpTest)

		require.Equal(t, 5, len(result))
		require.Equal(t, 1, result[0].Id)
		require.Equal(t, "test_name_1", result[0].Name)
		require.Equal(t, 23, result[0].Number)
		require.NotNil(t, result[0].CreatedAt)
		require.NotNil(t, result[0].UpdatedAt)
		require.Nil(t, result[0].DeletedAt)

		row, err = query.Use(c).WithLimit(5).
			WithOffset(5).
			All()
		require.NoError(t, err)
		require.NotNil(t, row)

		result = row.([]*CpTest)

		require.Equal(t, 5, len(result))
		require.Equal(t, 6, result[0].Id)
		require.Equal(t, "test_name_6", result[0].Name)
		require.Equal(t, 36, result[0].Number)
		require.NotNil(t, result[0].CreatedAt)
		require.Nil(t, result[0].UpdatedAt)
		require.Nil(t, result[0].DeletedAt)
	})

	t.Run("Aggregate", func(t *testing.T) {
		row, err := query.Use(c).
			WhereOp("number", OP_LESS_THAN, 25).
			Count()
		require.NoError(t, err)
		require.Equal(t, float64(5), *row)

		row, err = query.Use(c).
			WhereNil(DELETED_AT).
			Avg("number")
		require.NoError(t, err)
		require.Equal(t, float64(241)/10, *row)

		row, err = query.Use(c).
			AvgWithTrash("number")
		require.NoError(t, err)
		require.Equal(t, float64(280)/12, *row)

		row, err = query.Use(c).
			WhereNotNil(CREATED_AT).
			Sum("number")
		require.NoError(t, err)
		require.Equal(t, float64(241), *row)

		row, err = query.Use(c).
			SumWithTrash("number")
		require.NoError(t, err)
		require.Equal(t, float64(280), *row)

		row, err = query.Use(c).
			Min("number")
		require.NoError(t, err)
		require.Equal(t, float64(5), *row)

		row, err = query.Use(c).
			MinWithTrash("number")
		require.NoError(t, err)
		require.Equal(t, float64(5), *row)

		row, err = query.Use(c).
			Max("number")
		require.NoError(t, err)
		require.Equal(t, float64(40), *row)

		row, err = query.Use(c).
			MaxWithTrash("number")
		require.NoError(t, err)
		require.Equal(t, float64(40), *row)
	})

	t.Run("Bulk Insert", func(t *testing.T) {
		var cs []*CpTest
		for i := 0; i < 3; i++ {
			cp := NewCpTest()
			cp.Name = fmt.Sprintf("bulk_%d", i)
			cp.Number = i * i

			cs = append(cs, cp)
		}

		success, err := query.BulkInsert(cs)
		require.NoError(t, err)
		require.True(t, success)

		cnt, err := query.Use(c).CountWithTrash()
		require.NoError(t, err)
		require.Equal(t, float64(15), *cnt)
	})
}

func TestQuery_Raw(t *testing.T) {
	type RawTest struct {
		Name   string `db:"name"`
		Number int    `db:"number"`
	}

	sample := &RawTest{}

	query := initQuery()
	initTestDb(query)
	initTestData(query)

	c = NewCpTest()

	t.Run("Raw Named", func(t *testing.T) {
		payload := map[string]interface{}{
			"number": 5,
		}

		row, err := query.RawNamedAll(sample, "SELECT name, number FROM complete_tests WHERE number > :number", payload)
		require.NoError(t, err)
		require.NotNil(t, row)

		results := row.([]*RawTest)
		require.Equal(t, 11, len(results))
		require.Equal(t, "test_name_12", results[10].Name)
		require.Equal(t, 18, results[10].Number)

		row, err = query.RawNamedFirst(sample, "SELECT name, number FROM complete_tests WHERE number = :number", payload)
		require.NoError(t, err)
		require.NotNil(t, row)

		result := row.(*RawTest)
		require.Equal(t, "test_name_2", result.Name)
		require.Equal(t, 5, result.Number)

		insertPayload := map[string]interface{}{
			"name":   "insert_exec_test",
			"number": 1,
		}

		success, err := query.RawNamedExec("INSERT INTO complete_tests (name, number) VALUES (:name, :number);", insertPayload)
		require.NoError(t, err)
		require.True(t, success)

		row, err = query.RawNamedFirst(sample, "SELECT name, number FROM complete_tests WHERE id = 13;", nil)
		require.NoError(t, err)
		require.NotNil(t, row)

		result = row.(*RawTest)
		require.Equal(t, "insert_exec_test", result.Name)
		require.Equal(t, 1, result.Number)
	})

	t.Run("Raw", func(t *testing.T) {
		payload := []interface{}{
			5,
		}

		row, err := query.RawAll(sample, "SELECT name, number FROM complete_tests WHERE number <= $1 ORDER BY number ASC;", payload...)
		require.NoError(t, err)
		require.NotNil(t, row)

		results := row.([]*RawTest)
		require.Equal(t, 2, len(results))
		require.Equal(t, "insert_exec_test", results[0].Name)
		require.Equal(t, 1, results[0].Number)

		deletePayload := []interface{}{
			13,
		}

		success, err := query.RawExec("DELETE FROM complete_tests WHERE id = $1;", deletePayload...)
		require.NoError(t, err)
		require.True(t, success)

		firstPayload := []interface{}{
			1,
		}

		row, err = query.RawFirst(sample, "SELECT name, number FROM complete_tests WHERE number = $1;", firstPayload...)
		require.Error(t, err)
		require.Nil(t, row)
	})
}

func TestQuery_GroupBy(t *testing.T) {
	type GroupTest struct {
		Model
		Id     int    `db:"int"`
		Name   string `db:"name"`
		Number int    `db:"int"`
		Sum    int    `db:"sum"`
	}

	sample := &GroupTest{
		Model: AutoIncrementModel(
			"complete_tests",
			[]Column{
				Col("id"),
				Col("name"),
				Col("number"),
			},
			"id",
			true,
			true,
		),
	}

	query := initQuery()
	initTestDb(query)
	initGroupByData(query)

	c = NewCpTest()

	row, err := query.Use(sample).
		Columns(Col("name"), Sum("number")).
		GroupBy("name").
		All()
	require.NoError(t, err)
	require.NotNil(t, row)

	result := row.([]*GroupTest)

	require.Equal(t, 4, len(result))
	require.Equal(t, "test_name_1", result[0].Name)
	require.Equal(t, 47, result[0].Sum)

	row, err = query.Use(sample).
		Columns(Col("name"), Sum("number")).
		GroupBy("name").
		HavingOp(Sum("number"), OP_LESS_THAN, 50).
		OrHavingOp(Sum("number"), OP_MORE_THAN, 200).
		OrderByAggregate(Sum("number"), DIR_ASC).
		All()
	require.NoError(t, err)
	require.NotNil(t, row)

	result = row.([]*GroupTest)

	require.Equal(t, 2, len(result))
	require.Equal(t, "test_name_4", result[0].Name)
	require.Equal(t, 31, result[0].Sum)

	row, err = query.Use(sample).
		Columns(Col("name"), Sum("number")).
		GroupBy("name").
		Having(Sum("number"), 47).
		OrderByAggregate(Sum("number"), DIR_ASC).
		All()
	require.NoError(t, err)
	require.NotNil(t, row)

	result = row.([]*GroupTest)

	require.Equal(t, 1, len(result))
	require.Equal(t, "test_name_1", result[0].Name)
	require.Equal(t, 47, result[0].Sum)

	row, err = query.Use(sample).
		Columns(Col("name"), Sum("number")).
		GroupBy("name").
		GroupHaving(func(query *Query) *Query {
			return query.Having(Sum("number"), 107).
				OrHaving(Sum("number"), 47).
				OrHaving(Sum("number"), 95)
		}).
		OrderByAggregate(Sum("number"), DIR_DESC).
		All()
	require.NoError(t, err)
	require.NotNil(t, row)

	result = row.([]*GroupTest)

	require.Equal(t, 3, len(result))
	require.Equal(t, "test_name_2", result[0].Name)
	require.Equal(t, 107, result[0].Sum)
}
