package fazzdb

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"reflect"
)

// QueryDb creates a new pointer to the Query instance using a sqlx.DB instance and config struct,
// this constructor will automatically create a new transaction for query used in this Query instance
// this Query instance will automatically commit and rollback on the query it runs
func QueryDb(db *sqlx.DB, config Config) (*Query, error) {
	tx, err := db.Beginx()
	if nil != err {
		return nil, err
	}

	q := QueryTx(tx, config)
	q.AutoCommit = true

	return q, nil
}

// QueryTx creates a new pointer to the Query instance using a sqlx.Tx instance and config struct,
// this constructor will use the provided transaction and will not commit or rollback any query it runs
func QueryTx(tx *sqlx.Tx, config Config) *Query {
	return &Query{
		Config:     config,
		Parameter:  NewParameter(config),
		Model:      nil,
		Builder:    &Builder{},
		Tx:         tx,
		AutoCommit: false,
	}
}

// Query is a struct that will handle query building and struct mapping to the database
type Query struct {
	*Parameter
	Config     Config
	Model      ModelInterface
	Builder    *Builder
	Tx         *sqlx.Tx
	AutoCommit bool
}

// RawExec is a function that will run exec to a raw query with provided payload
func (q *Query) RawExec(query string, payload ...interface{}) (bool, error) {
	_, err := q.Tx.Exec(query, payload)
	if nil != err {
		q.autoRollback()
		return false, err
	}
	return true, err
}

// RawFirst is a function that will run raw query that return only one result with provided payload
func (q *Query) RawFirst(sample interface{}, query string, payload ...interface{}) (interface{}, error) {
	result, err := q.makeTypeOf(sample)
	if nil != err {
		return nil, err
	}

	stmt, err := q.Tx.Preparex(query)
	if nil != err {
		return nil, err
	}

	err = stmt.Get(result, payload...)
	if nil != err {
		return nil, err
	}

	return reflect.ValueOf(result).Elem().Interface(), nil
}

// RawAll is a function that will run raw query that return multiple result with provided payload
func (q *Query) RawAll(sample interface{}, query string, payload ...interface{}) (interface{}, error) {
	results, err := q.makeSliceOf(sample)
	if nil != err {
		return nil, err
	}

	stmt, err := q.Tx.Preparex(query)
	if nil != err {
		return nil, err
	}

	err = stmt.Select(results, payload...)
	if nil != err {
		return nil, err
	}

	return reflect.ValueOf(results).Elem().Interface(), nil
}

// RawNamedExec is a function that will run exec to a raw named query with provided payload
func (q *Query) RawNamedExec(query string, payload map[string]interface{}) (bool, error) {
	stmt, err := q.Tx.PrepareNamed(query)
	if nil != err {
		q.autoRollback()
		return false, err
	}

	_, err = stmt.Exec(payload)
	if nil != err {
		q.autoRollback()
		return false, err
	}

	q.autoCommit()
	return true, nil
}

// RawNamedFirst is a function that will run raw named query that return only one result with provided payload
func (q *Query) RawNamedFirst(sample interface{}, query string, payload map[string]interface{}) (interface{}, error) {
	result, err := q.makeTypeOf(sample)
	if nil != err {
		return nil, err
	}

	stmt, err := q.Tx.PrepareNamed(query)
	if nil != err {
		return nil, err
	}

	err = stmt.Get(result, payload)
	if nil != err {
		return nil, err
	}

	return reflect.ValueOf(result).Elem().Interface(), nil
}

// RawNamedAll is a function that will run raw named query that return multiple result with provided payload
func (q *Query) RawNamedAll(sample interface{}, query string, payload map[string]interface{}) (interface{}, error) {
	results, err := q.makeSliceOf(sample)
	if nil != err {
		return nil, err
	}

	stmt, err := q.Tx.PrepareNamed(query)
	if nil != err {
		return nil, err
	}

	err = stmt.Select(results, payload)
	if nil != err {
		return nil, err
	}

	return reflect.ValueOf(results).Elem().Interface(), nil
}

// First is a function that will return query with only one result
func (q *Query) First() (interface{}, error) {
	return q.first(NO_TRASH)
}

// FirstWithTrash is a function that will return query with only one result including soft deleted row
func (q *Query) FirstWithTrash() (interface{}, error) {
	return q.first(WITH_TRASH)
}

// All is a function that will return query with multiple result
func (q *Query) All() (interface{}, error) {
	return q.all(NO_TRASH)
}

// AllWithTrash is a function that will return query with multiple result including soft deleted row
func (q *Query) AllWithTrash() (interface{}, error) {
	return q.all(WITH_TRASH)
}

// Insert is a function that will insert data based on model attribute
func (q *Query) Insert() (*interface{}, error) {
	var id interface{}

	err := q.handleNilModel()
	if nil != err {
		return nil, err
	}

	q.Model.GeneratePK()

	if q.Model.IsTimestamps() {
		q.Model.created()
	}

	query := q.Builder.BuildInsert(q.Model)

	stmt, err := q.Tx.PrepareNamed(query)
	if nil != err {
		q.autoRollback()
		return nil, err
	}

	err = stmt.Get(&id, q.Model.Payload())
	if nil != err {
		q.autoRollback()
		return nil, err
	}

	q.autoCommit()
	return &id, nil
}

// BulkInsert is a function that will insert multiple data in one query, receive slice of model
func (q *Query) BulkInsert(data interface{}) (bool, error) {
	err := q.handleNilModel()
	if nil != err {
		return false, err
	}

	d := reflect.ValueOf(data)
	if d.Kind() != reflect.Slice {
		return false, fmt.Errorf("payload must be a slice")
	}

	slice := make([]interface{}, d.Len())
	for i := 0; i < d.Len(); i++ {
		slice[i] = d.Index(i).Interface()
	}

	query := q.Builder.BuildBulkInsert(q.Model, slice)
	payloads := q.bulkPayload(slice)

	stmt, err := q.Tx.PrepareNamed(query)
	if nil != err {
		q.autoRollback()
		return false, err
	}

	_, err = stmt.Exec(payloads)
	if nil != err {
		q.autoRollback()
		return false, err
	}

	q.autoCommit()
	return true, nil
}

// Update is a function that will update data based on model attribute with primary key attribute
func (q *Query) Update() (bool, error) {
	defer q.clearParameter()

	err := q.handleNilModel()
	if nil != err {
		return false, err
	}

	q.setPKCondition()

	if q.Model.IsTimestamps() {
		q.Model.updated()
	}

	query := q.Builder.BuildUpdate(q.Model, q.Parameter)
	query = q.bindIn(query)

	stmt, err := q.Tx.PrepareNamed(query)
	if nil != err {
		q.autoRollback()
		return false, err
	}

	_, err = stmt.Exec(q.mergedPayload())
	if nil != err {
		q.autoRollback()
		return false, err
	}

	q.autoCommit()
	return true, nil
}

// Delete is a function that will delete data based on model attribute with primary key attribute
// will automatically soft delete if soft delete attribute is active
func (q *Query) Delete() (bool, error) {
	defer q.clearParameter()

	err := q.handleNilModel()
	if nil != err {
		return false, err
	}

	q.setPKCondition()

	if q.Model.IsSoftDelete() {
		q.Model.deleted()
		return q.Update()
	}

	query := q.Builder.BuildDelete(q.Model, q.Parameter)
	query = q.bindIn(query)

	stmt, err := q.Tx.PrepareNamed(query)
	if nil != err {
		q.autoRollback()
		return false, err
	}

	_, err = stmt.Exec(q.mergedPayload())
	if nil != err {
		q.autoRollback()
		return false, err
	}

	q.autoCommit()
	return true, nil
}

func (q *Query) Aggregate(aggregate Aggregate, column string) (*float64, error) {
	return q.aggregate(aggregate, column, NO_TRASH)
}

func (q *Query) AggregateWithTrash(aggregate Aggregate, column string) (*float64, error) {
	return q.aggregate(aggregate, column, WITH_TRASH)
}

// Avg is a function that will return average of a column
func (q *Query) Avg(column string) (*float64, error) {
	return q.Aggregate(AG_AVG, column)
}

// AvgWithTrash is a function that will return average of a column with soft deleted row
func (q *Query) AvgWithTrash(column string) (*float64, error) {
	return q.AggregateWithTrash(AG_AVG, column)
}

// Min is a function that will return minimum of a column
func (q *Query) Min(column string) (*float64, error) {
	return q.Aggregate(AG_MIN, column)
}

// MinWithTrash is a function that will return minimum of a column with soft deleted row
func (q *Query) MinWithTrash(column string) (*float64, error) {
	return q.AggregateWithTrash(AG_MIN, column)
}

// Max is a function that will return maximum of a column
func (q *Query) Max(column string) (*float64, error) {
	return q.Aggregate(AG_MAX, column)
}

// MaxWithTrash is a function that will return maximum of a column with soft deleted row
func (q *Query) MaxWithTrash(column string) (*float64, error) {
	return q.AggregateWithTrash(AG_MAX, column)
}

// Sum is a function that will return sum of a column
func (q *Query) Sum(column string) (*float64, error) {
	return q.Aggregate(AG_SUM, column)
}

// SumWithTrash is a function that will return sum of a column with soft deleted row
func (q *Query) SumWithTrash(column string) (*float64, error) {
	return q.AggregateWithTrash(AG_SUM, column)
}

// Count is a function that will return count of a column
func (q *Query) Count() (*float64, error) {
	return q.Aggregate(AG_COUNT, "*")
}

// MinWithTrash is a function that will return count of a column with soft deleted row
func (q *Query) CountWithTrash() (*float64, error) {
	return q.AggregateWithTrash(AG_COUNT, "*")
}

// Use is a function that will set Model instance that will be used for query
func (q *Query) Use(m ModelInterface) *Query {
	q.Model = m
	return q
}