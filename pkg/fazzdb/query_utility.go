package fazzdb

import (
	"github.com/jmoiron/sqlx"
	"log"
	"reflect"
)

func (q *Query) assignModelSlices(results interface{}, m *Model) interface{} {
	slice := reflect.ValueOf(results).Elem().Interface()
	sVal := reflect.ValueOf(slice)
	for i := 0; i < sVal.Len(); i++ {
		assigned := q.assignModel(sVal.Index(i).Addr().Interface(), m)
		sVal.Index(i).Set(reflect.ValueOf(assigned))
	}
	return sVal.Interface()
}

func (q *Query) assignModel(result interface{}, m *Model) interface{} {
	value := reflect.ValueOf(result).Interface()
	model := reflect.ValueOf(&m).Elem()
	complete := reflect.ValueOf(value).Elem()
	complete.FieldByName("Model").Set(model)
	return complete.Interface()
}

func (q *Query) prepareSelect() (*sqlx.NamedStmt, map[string]interface{}, error) {
	query := q.Builder.BuildSelect(q.Model, q.Parameter)
	log.Println(query)
	stmt, err := q.Tx.PrepareNamed(query)
	if nil != err {
		return nil, nil, err
	}

	return stmt, q.Parameter.Values, err
}

func (q *Query) setPrimaryKeyParameter() {
	pkConditionExist := false
	for _, condition := range q.Conditions {
		if condition.Key == q.Model.GetPK() {
			pkConditionExist = true
			break
		}
	}

	if !pkConditionExist {
		q.Where(q.Model.GetPK(), q.Model.Get(q.Model.GetPK()))
	}
}

func (q *Query) mergedPayload() map[string]interface{} {
	payload := q.Model.Payload()
	for i, v := range q.Parameter.Values {
		payload[i] = v
	}
	return payload
}

func (q *Query) autoCommit() {
	if q.AutoCommit {
		_ = q.Tx.Commit()
	}
}

func (q *Query) autoRollback() {
	if q.AutoCommit {
		_ = q.Tx.Rollback()
	}
}

func (q *Query) clearParameter() {
	q.Parameter = NewParameter()
}