package fazzdb

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestModel_ColumnCount(t *testing.T) {
	initDefaultCpTest()

	require.Equal(t, 3, ai.ColumnCount())
	require.Equal(t, 3, u.ColumnCount())
	require.Equal(t, 3, p.ColumnCount())
	require.Equal(t, 3, ts.ColumnCount())
	require.Equal(t, 2, sd.ColumnCount())
	require.Equal(t, 6, c.ColumnCount())
}

func TestModel_GetTable(t *testing.T) {
	initDefaultCpTest()

	require.Equal(t, "auto_increment_tests", ai.GetTable())
	require.Equal(t, "uuid_tests", u.GetTable())
	require.Equal(t, "plain_tests", p.GetTable())
	require.Equal(t, "timestamp_tests", ts.GetTable())
	require.Equal(t, "soft_delete_tests", sd.GetTable())
	require.Equal(t, "complete_tests", c.GetTable())
}

func TestModel_GetColumns(t *testing.T) {
	initDefaultCpTest()

	require.Equal(t, []Column{Col("id"), Col("name"), Col("number")}, ai.GetColumns())
	require.Equal(t, []Column{Col("id"), Col("name"), Col("number")}, u.GetColumns())
	require.Equal(t, []Column{Col("id"), Col("name"), Col("number")}, p.GetColumns())
	require.Equal(t, []Column{Col("id"), Col("createdAt"), Col("updatedAt")}, ts.GetColumns())
	require.Equal(t, []Column{Col("id"), Col("deletedAt")}, sd.GetColumns())
	require.Equal(t, []Column{Col("id"), Col("name"), Col("number"), Col("createdAt"),
		Col("updatedAt"), Col("deletedAt")}, c.GetColumns())
}

func TestModel_IsAutoIncrement(t *testing.T) {
	initDefaultCpTest()

	require.Equal(t, true, ai.IsAutoIncrement())
	require.Equal(t, false, u.IsAutoIncrement())
	require.Equal(t, false, p.IsAutoIncrement())
	require.Equal(t, true, ts.IsAutoIncrement())
	require.Equal(t, true, sd.IsAutoIncrement())
	require.Equal(t, true, c.IsAutoIncrement())
}

func TestModel_IsTimestamps(t *testing.T) {
	initDefaultCpTest()

	require.Equal(t, false, ai.IsTimestamps())
	require.Equal(t, false, u.IsTimestamps())
	require.Equal(t, false, p.IsTimestamps())
	require.Equal(t, true, ts.IsTimestamps())
	require.Equal(t, false, sd.IsTimestamps())
	require.Equal(t, true, c.IsTimestamps())
}

func TestModel_IsSoftDelete(t *testing.T) {
	initDefaultCpTest()

	require.Equal(t, false, ai.IsSoftDelete())
	require.Equal(t, false, u.IsSoftDelete())
	require.Equal(t, false, p.IsSoftDelete())
	require.Equal(t, false, ts.IsSoftDelete())
	require.Equal(t, true, sd.IsSoftDelete())
	require.Equal(t, true, c.IsSoftDelete())
}

func TestModel_IsUuid(t *testing.T) {
	initDefaultCpTest()

	require.Equal(t, false, ai.IsUuid())
	require.Equal(t, true, u.IsUuid())
	require.Equal(t, false, p.IsUuid())
	require.Equal(t, false, ts.IsUuid())
	require.Equal(t, false, sd.IsUuid())
	require.Equal(t, false, c.IsUuid())
}

func TestModel_Get(t *testing.T) {
	initDefaultCpTest()

	assert.Panics(t, func() { ai.Get("name") })
	assert.Panics(t, func() { u.Get("name") })
	assert.Panics(t, func() { p.Get("name") })
	assert.Panics(t, func() { ts.Get("createdAt") })
	assert.Panics(t, func() { sd.Get("deletedAt") })
	require.Equal(t, "c_name_test", c.Get("name"))
}

func TestModel_GeneratePK(t *testing.T) {
	current := u.Id
	u.GeneratePK()
	require.NotEqual(t, current, u.Id)
}

func TestModel_GetModel(t *testing.T) {
	initDefaultCpTest()

	require.Equal(
		t,
		*AutoIncrementModel(
			"auto_increment_tests",
			[]Column{
				Col("id"),
				Col("name"),
				Col("number"),
			},
			"id",
			false,
			false,
		),
		ai.GetModel(),
	)

	require.Equal(
		t,
		*UuidModel(
			"uuid_tests",
			[]Column{
				Col("id"),
				Col("name"),
				Col("number"),
			},
			"id",
			false,
			false,
		),
		u.GetModel(),
	)

	require.Equal(
		t,
		*PlainModel(
			"plain_tests",
			[]Column{
				Col("id"),
				Col("name"),
				Col("number"),
			},
			"id",
			false,
			false,
		),
		p.GetModel(),
	)

	require.Equal(
		t,
		*AutoIncrementModel(
			"timestamp_tests",
			[]Column{
				Col("id"),
			},
			"id",
			true,
			false,
		),
		ts.GetModel(),
	)

	require.Equal(
		t,
		*AutoIncrementModel(
			"soft_delete_tests",
			[]Column{
				Col("id"),
			},
			"id",
			false,
			true,
		),
		sd.GetModel(),
	)

	require.Equal(
		t,
		*AutoIncrementModel(
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
		c.GetModel(),
	)
}