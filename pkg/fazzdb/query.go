package fazzdb

import (
	"github.com/jmoiron/sqlx"
	"log"
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

	log.Println(stmt)

	err = stmt.Select(results, payload)
	if nil != err {
		return nil, err
	}

	return reflect.ValueOf(results).Elem().Interface(), nil
}

func (q *Query) First() (interface{}, error) {
	defer q.clearParameter()

	result, err := q.makeTypeOf(q.Model)
	if nil != err {
		return nil, err
	}

	q.setLimit(1)
	stmt, args, err := q.prepareSelect(AG_NONE, "")
	if nil != err {
		return nil, err
	}

	err = stmt.Get(result, args)
	if nil != err {
		return nil, err
	}

	return q.assignModel(result, q.Model.GetModel()), nil
}

func (q *Query) All() (interface{}, error) {
	defer q.clearParameter()

	results, err := q.makeSliceOf(q.Model)
	if nil != err {
		return nil, err
	}

	stmt, args, err := q.prepareSelect(AG_NONE, "")
	if nil != err {
		return nil, err
	}

	err = stmt.Select(results, args)
	if nil != err {
		return nil, err
	}

	return q.assignModelSlices(results, q.Model.GetModel()), nil
}

func (q *Query) Insert() (*interface{}, error) {
	var id interface{}

	q.Model.GeneratePK()

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

func (q *Query) Update() (bool, error) {
	defer q.clearParameter()

	q.setPKCondition()

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

	q.setPKCondition()

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
	defer q.clearParameter()

	var result float64

	stmt, args, err := q.prepareSelect(aggregate, column)
	if nil != err {
		return nil, err
	}

	err = stmt.Get(&result, args)
	if nil != err {
		return nil, err
	}

	return &result, nil
}

func (q *Query) Avg(column string) (*float64, error) {
	return q.Aggregate(AG_AVG, column)
}

func (q *Query) Min(column string) (*float64, error) {
	return q.Aggregate(AG_MIN, column)
}

func (q *Query) Max(column string) (*float64, error) {
	return q.Aggregate(AG_MAX, column)
}

func (q *Query) Sum(column string) (*float64, error) {
	return q.Aggregate(AG_SUM, column)
}

func (q *Query) Count() (*float64, error) {
	return q.Aggregate(AG_COUNT, "*")
}

func (q *Query) Use(m ModelInterface) *Query {
	q.Model = m
	return q
}