package fazzdb

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"reflect"
)

func QueryDb(db *sqlx.DB, config Config) (*Query, error) {
	tx, err := db.Beginx()
	if nil != err {
		return nil, err
	}

	q := QueryTx(tx, config)
	q.AutoCommit = true

	return q, nil
}

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

type Query struct {
	*Parameter
	Config     Config
	Model      ModelInterface
	Builder    *Builder
	Tx         *sqlx.Tx
	AutoCommit bool
}

func (q *Query) RawExec(query string, payload ...interface{}) (bool, error) {
	_, err := q.Tx.Exec(query, payload)
	if nil != err {
		q.autoRollback()
		return false, err
	}
	return true, err
}

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

func (q *Query) First() (interface{}, error) {
	return q.first(NO_TRASH)
}

func (q *Query) FirstWithTrash() (interface{}, error) {
	return q.first(WITH_TRASH)
}

func (q *Query) All() (interface{}, error) {
	return q.all(NO_TRASH)
}

func (q *Query) AllWithTrash() (interface{}, error) {
	return q.all(WITH_TRASH)
}

func (q *Query) Insert() (*interface{}, error) {
	var id interface{}

	err := q.handleNilModel()
	if nil != err {
		return nil, err
	}

	q.Model.GeneratePK()

	if q.Model.IsTimestamps() {
		q.Model.Created()
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

func (q *Query) Update() (bool, error) {
	defer q.clearParameter()

	err := q.handleNilModel()
	if nil != err {
		return false, err
	}

	q.setPKCondition()

	if q.Model.IsTimestamps() {
		q.Model.Updated()
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

func (q *Query) Delete() (bool, error) {
	defer q.clearParameter()

	err := q.handleNilModel()
	if nil != err {
		return false, err
	}

	q.setPKCondition()

	if q.Model.IsSoftDelete() {
		q.Model.Deleted()
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

func (q *Query) Avg(column string) (*float64, error) {
	return q.Aggregate(AG_AVG, column)
}

func (q *Query) AvgWithTrash(column string) (*float64, error) {
	return q.AggregateWithTrash(AG_AVG, column)
}

func (q *Query) Min(column string) (*float64, error) {
	return q.Aggregate(AG_MIN, column)
}

func (q *Query) MinWithTrash(column string) (*float64, error) {
	return q.AggregateWithTrash(AG_MIN, column)
}

func (q *Query) Max(column string) (*float64, error) {
	return q.Aggregate(AG_MAX, column)
}

func (q *Query) MaxWithTrash(column string) (*float64, error) {
	return q.AggregateWithTrash(AG_MAX, column)
}

func (q *Query) Sum(column string) (*float64, error) {
	return q.Aggregate(AG_SUM, column)
}

func (q *Query) SumWithTrash(column string) (*float64, error) {
	return q.AggregateWithTrash(AG_SUM, column)
}

func (q *Query) Count() (*float64, error) {
	return q.Aggregate(AG_COUNT, "*")
}

func (q *Query) CountWithTrash() (*float64, error) {
	return q.AggregateWithTrash(AG_COUNT, "*")
}

func (q *Query) Use(m ModelInterface) *Query {
	q.Model = m
	return q
}