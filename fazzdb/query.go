package fazzdb

import (
	"db/fazzdb/fazzorder"
	"db/fazzdb/fazzspec"
	"github.com/jmoiron/sqlx"
	"log"
	"reflect"
)

// TODO: error handling with stack trace

type Query struct {
	Parameter
	Model      ModelInterface
	Builder    *Builder
	Tx         *sqlx.Tx
	AutoCommit bool
}

func (q *Query) First() (interface{}, error) {
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

	log.Println(reflect.ValueOf(results).Elem())

	return reflect.ValueOf(results).Elem().Interface(), nil
}

func (q *Query) Insert() (*uint64, error) {
	var id uint64

	query := q.Builder.BuildInsert(q.Model)
	stmt, err := q.Tx.PrepareNamed(query)
	if nil != err {
		return nil, err
	}

	err = stmt.Get(&id, q.Model.Payload())
	if nil != err {
		return nil, err
	}

	return &id, nil
}

func (q *Query) Update() (bool, error) {
	pkConditionExist := false
	for _, condition := range q.Conditions {
		if condition.Key == q.Model.GetPrimaryKey() {
			pkConditionExist = true
			break
		}
	}

	if !pkConditionExist {
		q.Where(q.Model.GetPrimaryKey(), q.Model.Get(q.Model.GetPrimaryKey()))
	}

	log.Println(q.Model.Payload())
	query := q.Builder.BuildUpdate(q.Model, &q.Parameter)
	log.Println(query)
	stmt, err := q.Tx.PrepareNamed(query)
	if nil != err {
		panic(err)
		return false, err
	}

	_, err = stmt.Exec(q.Model.Payload())
	if nil != err {
		panic(err)
		return false, err
	}

	return true, nil
}

func (q *Query) Where(key string, value interface{}) *Query {
	return q.AppendCondition(fazzspec.CO_AND, key, fazzspec.OP_EQUALS, value)
}

func (q *Query) WhereWithOperator(key string, operator fazzspec.Operator, value interface{}) *Query {
	return q.AppendCondition(fazzspec.CO_AND, key, operator, value)
}

func (q *Query) OrWhere(key string, value interface{}) *Query {
	return q.AppendCondition(fazzspec.CO_OR, key, fazzspec.OP_EQUALS, value)
}

func (q *Query) OrWhereWithOperator(key string, operator fazzspec.Operator, value interface{}) *Query {
	return q.AppendCondition(fazzspec.CO_OR, key, operator, value)
}

func (q *Query) OrderBy(key string, direction fazzorder.OrderDirection) *Query {
	q.appendOrder(q.Model.GetTable(), key, direction)
	return q
}

func (q *Query) WithLimit(limit int) *Query {
	q.setLimit(limit)
	return q
}

func (q *Query) WithOffset(offset int) *Query {
	q.setOffset(offset)
	return q
}

func (q *Query) WithLock(lock fazzspec.Lock) *Query {
	q.setLock(lock)
	return q
}

func (q *Query) AppendCondition(connector fazzspec.Connector, key string, operator fazzspec.Operator, value interface{}) *Query {
	q.appendCondition(q.Model.GetTable(), connector, key, operator, value)
	return q
}

func (q *Query) assignModel(result interface{}, m *Model) interface{} {
	complete := reflect.ValueOf(result).Interface()
	model := reflect.ValueOf(m).Elem()
	reflect.ValueOf(complete).Elem().FieldByName("Model").Set(model)
	return complete
}

func (q *Query) prepareSelect() (*sqlx.NamedStmt, map[string]interface{}, error) {
	query := q.Builder.BuildSelect(q.Model, &q.Parameter)
	stmt, err := q.Tx.PrepareNamed(query)
	if nil != err {
		return nil, nil, err
	}

	return stmt, q.Parameter.Values, err
}

func (q *Query) useDefaultIfNil(value interface{}, defaultValue interface{}) interface{} {
	// TODO: NEED TO CHECK IF Ptr or not
	if nil != value {
		return value
	}
	return defaultValue
}

func NewQuery(db *sqlx.DB, m ModelInterface) *Query {
	tx, err := db.Beginx()
	if nil != err {
		log.Fatal(err)
	}

	q := NewQueryTx(tx, m)
	q.AutoCommit = true

	return q
}

func NewQueryTx(tx *sqlx.Tx, m ModelInterface) *Query {
	return &Query{
		Parameter:  *NewParameter(),
		Model:      m,
		Builder:    &Builder{},
		Tx:         tx,
		AutoCommit: false,
	}
}
