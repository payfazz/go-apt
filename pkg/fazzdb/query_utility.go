package fazzdb

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"reflect"
	"strings"
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

func (q *Query) first(withTrash TrashStatus) (interface{}, error) {
	defer q.clearParameter()

	err := q.handleNilModel()
	if nil != err {
		return nil, err
	}

	result, err := q.makeTypeOf(q.Model)
	if nil != err {
		return nil, err
	}

	q.setLimit(1)
	stmt, args, err := q.prepareSelect(AG_NONE, "", withTrash)
	if nil != err {
		return nil, err
	}

	err = stmt.Get(result, args)
	if nil != err {
		return nil, err
	}

	return q.assignModel(result, q.Model.GetModel()), nil
}

func (q *Query) all(withTrash TrashStatus) (interface{}, error) {
	defer q.clearParameter()

	err := q.handleNilModel()
	if nil != err {
		return nil, err
	}

	results, err := q.makeSliceOf(q.Model)
	if nil != err {
		return nil, err
	}

	stmt, args, err := q.prepareSelect(AG_NONE, "", withTrash)
	if nil != err {
		return nil, err
	}

	err = stmt.Select(results, args)
	if nil != err {
		return nil, err
	}

	return q.assignModelSlices(results, q.Model.GetModel()), nil
}

func (q *Query) aggregate(aggregate Aggregate, column string, withTrash TrashStatus) (*float64, error) {
	defer q.clearParameter()

	err := q.handleNilModel()
	if nil != err {
		return nil, err
	}

	var result float64

	stmt, args, err := q.prepareSelect(aggregate, column, NO_TRASH)
	if nil != err {
		return nil, err
	}

	err = stmt.Get(&result, args)
	if nil != err {
		return nil, err
	}

	return &result, nil
}

func (q *Query) prepareSelect(aggregate Aggregate, aggregateColumn string, withTrash TrashStatus) (*sqlx.NamedStmt, map[string]interface{}, error) {
	if q.Model.IsSoftDelete() && withTrash == NO_TRASH {
		q.WhereNil("deletedAt")
	}

	if len(q.Parameter.Orders) == 0 {
		q.OrderBy(q.Model.GetPK(), DIR_ASC)
	}

	query := q.Builder.BuildSelect(q.Model, q.Parameter, aggregate, aggregateColumn)
	query = q.bindIn(query)

	stmt, err := q.Tx.PrepareNamed(query)
	if nil != err {
		return nil, nil, err
	}

	return stmt, q.Parameter.Values, err
}

func (q *Query) bindIn(query string) string {
	for i, value := range q.Parameter.Values {
		if reflect.TypeOf(value).Kind() == reflect.Slice {
			inValueQuery := ""
			sliceValue := reflect.ValueOf(value)
			for j := 0; j < sliceValue.Len(); j++ {
				prefix := fmt.Sprintf("%s%d", i, j)
				if j == 0 {
					inValueQuery = fmt.Sprintf("%s :%s", inValueQuery, prefix)
				} else {
					inValueQuery = fmt.Sprintf("%s, :%s", inValueQuery, prefix)
				}
				q.Parameter.Values[prefix] = sliceValue.Index(j).Interface()
			}
			query = strings.ReplaceAll(query, fmt.Sprintf(":%s", i), inValueQuery)
		}
	}

	return query
}

func (q *Query) setPKCondition() {
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

func (q *Query) bulkPayload(data []interface{}) map[string]interface{} {
	payloads := map[string]interface{}{}
	for i, v := range data {
		model := v.(ModelInterface)
		model.GeneratePK()
		payload := model.Payload()
		for key, value := range payload {
			payloads[fmt.Sprintf("%d%s", i, key)] = value
		}
	}
	return payloads
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

func (q *Query) handleNilModel() error {
	if nil == q.Model {
		return fmt.Errorf("please use a model before doing query")
	}
	return nil
}

func (q *Query) clearParameter() {
	q.Parameter = NewParameter(q.Config)
}

func (q *Query) makeTypeOf(sample interface{}) (interface{}, error) {
	if reflect.TypeOf(sample).Kind() != reflect.Ptr {
		return nil, fmt.Errorf("sample must be a pointer to reference model")
	}
	element := reflect.TypeOf(sample).Elem()
	return reflect.New(element).Interface(), nil
}

func (q *Query) makeSliceOf(sample interface{}) (interface{}, error) {
	if reflect.TypeOf(sample).Kind() != reflect.Ptr {
		return nil, fmt.Errorf("sample must be a pointer to reference model")
	}
	element := reflect.TypeOf(sample).Elem()
	return reflect.New(reflect.SliceOf(element)).Interface(), nil
}