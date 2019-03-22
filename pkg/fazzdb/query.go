package fazzdb

import (
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

func (q *Query) First() (interface{}, error) {
	defer q.clearParameter()

	element := reflect.TypeOf(q.Model).Elem()
	result := reflect.New(element).Interface()

	q.setLimit(1)
	stmt, args, err := q.prepareSelect()
	if nil != err {
		return nil, err
	}

	err = stmt.Get(result, args)
	if nil != err {
		return nil, err
	}

	return q.assignModel(result, q.Model.GetModel()), nil
}

func (q *Query) GetAll() (interface{}, error) {
	defer q.clearParameter()

	element := reflect.TypeOf(q.Model).Elem()
	results := reflect.New(reflect.SliceOf(element)).Interface()

	stmt, args, err := q.prepareSelect()
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

func (q *Query) Use(m ModelInterface) *Query {
	q.Model = m
	return q
}