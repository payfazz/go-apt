package fazzdb

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"reflect"
	"strings"
)

// QueryDb creates a new pointer to the Query instance using a sqlx.DB instance and config struct,
// this constructor will automatically create a new transaction for query used in this Query instance
// this Query instance will automatically commit and rollback on the query it runs
func QueryDb(db *sqlx.DB, config Config) *Query {
	return &Query{
		Config:     config,
		Parameter:  NewParameter(config),
		Model:      nil,
		Builder:    &Builder{},
		Db:         db,
		AutoCommit: true,
	}
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
	Db         *sqlx.DB
	Tx         *sqlx.Tx
	AutoCommit bool
}

// RawExec is a function that will run exec to a raw query with provided payload
func (q *Query) RawExec(query string, payload ...interface{}) (bool, error) {
	err := q.autoBegin()
	if nil != err {
		return false, err
	}

	_, err = q.Tx.Exec(query, payload)
	if nil != err {
		q.autoRollback()
		return false, err
	}

	q.autoCommit()
	return true, err
}

// RawFirst is a function that will run raw query that return only one result with provided payload
func (q *Query) RawFirst(sample interface{}, query string, payload ...interface{}) (interface{}, error) {
	err := q.autoBegin()
	if nil != err {
		return nil, err
	}

	result, err := q.makeTypeOf(sample)
	if nil != err {
		q.autoRollback()
		return nil, err
	}

	stmt, err := q.Tx.Preparex(query)
	if nil != err {
		q.autoRollback()
		return nil, err
	}

	err = stmt.Get(result, payload...)
	if nil != err {
		q.autoRollback()
		return nil, err
	}

	q.autoCommit()
	return reflect.ValueOf(result).Elem().Interface(), nil
}

// RawAll is a function that will run raw query that return multiple result with provided payload
func (q *Query) RawAll(sample interface{}, query string, payload ...interface{}) (interface{}, error) {
	err := q.autoBegin()
	if nil != err {
		return nil, err
	}

	results, err := q.makeSliceOf(sample)
	if nil != err {
		q.autoRollback()
		return nil, err
	}

	stmt, err := q.Tx.Preparex(query)
	if nil != err {
		q.autoRollback()
		return nil, err
	}

	err = stmt.Select(results, payload...)
	if nil != err {
		q.autoRollback()
		return nil, err
	}

	q.autoCommit()
	return reflect.ValueOf(results).Elem().Interface(), nil
}

// RawNamedExec is a function that will run exec to a raw named query with provided payload
func (q *Query) RawNamedExec(query string, payload map[string]interface{}) (bool, error) {
	err := q.autoBegin()
	if nil != err {
		return false, err
	}

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
	err := q.autoBegin()
	if nil != err {
		return nil, err
	}

	result, err := q.makeTypeOf(sample)
	if nil != err {
		q.autoRollback()
		return nil, err
	}

	stmt, err := q.Tx.PrepareNamed(query)
	if nil != err {
		q.autoRollback()
		return nil, err
	}

	err = stmt.Get(result, payload)
	if nil != err {
		q.autoRollback()
		return nil, err
	}

	q.autoCommit()
	return reflect.ValueOf(result).Elem().Interface(), nil
}

// RawNamedAll is a function that will run raw named query that return multiple result with provided payload
func (q *Query) RawNamedAll(sample interface{}, query string, payload map[string]interface{}) (interface{}, error) {
	err := q.autoBegin()
	if nil != err {
		return nil, err
	}

	results, err := q.makeSliceOf(sample)
	if nil != err {
		q.autoRollback()
		return nil, err
	}

	stmt, err := q.Tx.PrepareNamed(query)
	if nil != err {
		q.autoRollback()
		return nil, err
	}

	err = stmt.Select(results, payload)
	if nil != err {
		q.autoRollback()
		return nil, err
	}

	q.autoCommit()
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

	err = q.autoBegin()
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

	err = q.autoBegin()
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

	err = q.autoBegin()
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

	err = q.autoBegin()
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

// Avg is a function that will return average of a column
func (q *Query) Avg(column string) (*float64, error) {
	return q.aggregate(AG_AVG, column, NO_TRASH)
}

// AvgWithTrash is a function that will return average of a column with soft deleted row
func (q *Query) AvgWithTrash(column string) (*float64, error) {
	return q.aggregate(AG_AVG, column, WITH_TRASH)
}

// Min is a function that will return minimum of a column
func (q *Query) Min(column string) (*float64, error) {
	return q.aggregate(AG_MIN, column, NO_TRASH)
}

// MinWithTrash is a function that will return minimum of a column with soft deleted row
func (q *Query) MinWithTrash(column string) (*float64, error) {
	return q.aggregate(AG_MIN, column, WITH_TRASH)
}

// Max is a function that will return maximum of a column
func (q *Query) Max(column string) (*float64, error) {
	return q.aggregate(AG_MAX, column, NO_TRASH)
}

// MaxWithTrash is a function that will return maximum of a column with soft deleted row
func (q *Query) MaxWithTrash(column string) (*float64, error) {
	return q.aggregate(AG_MAX, column, WITH_TRASH)
}

// Sum is a function that will return sum of a column
func (q *Query) Sum(column string) (*float64, error) {
	return q.aggregate(AG_SUM, column, NO_TRASH)
}

// SumWithTrash is a function that will return sum of a column with soft deleted row
func (q *Query) SumWithTrash(column string) (*float64, error) {
	return q.aggregate(AG_SUM, column, WITH_TRASH)
}

// Count is a function that will return count of a column
func (q *Query) Count() (*float64, error) {
	return q.aggregate(AG_COUNT, "*", NO_TRASH)
}

// MinWithTrash is a function that will return count of a column with soft deleted row
func (q *Query) CountWithTrash() (*float64, error) {
	return q.aggregate(AG_COUNT, "*", WITH_TRASH)
}

// Use is a function that will set Model instance that will be used for query
func (q *Query) Use(m ModelInterface) *Query {
	q.Model = m
	return q
}

// Where is a function that will add new condition that will check if column equals value given, connector
// that is used between condition is AND connector
func (q *Query) Where(key string, value interface{}) *Query {
	return q.AppendCondition(CO_AND, key, OP_EQUALS, value)
}

// WhereOp is a function that will add new condition that will check if column fulfill operator with value given,
// connector that is used between condition is AND connector
func (q *Query) WhereOp(key string, operator Operator, value interface{}) *Query {
	return q.AppendCondition(CO_AND, key, operator, value)
}

// WhereNil is a function that will add new condition that will check if column nil, connector that is used
// between condition is AND connector
func (q *Query) WhereNil(key string) *Query {
	return q.AppendCondition(CO_AND, key, OP_IS_NULL, nil)
}

// WhereNotNil is a function that will add new condition that will check if column not nil, connector that is used
// between condition is AND connector
func (q *Query) WhereNotNil(key string) *Query {
	return q.AppendCondition(CO_AND, key, OP_IS_NOT_NULL, nil)
}

// OrWhere is a function that will add new condition that will check if column equals value given, connector
// that is used between condition is OR connector
func (q *Query) OrWhere(key string, value interface{}) *Query {
	return q.AppendCondition(CO_OR, key, OP_EQUALS, value)
}

// OrWhereOp is a function that will add new condition that will check if column fulfill operator with value given,
// connector that is used between condition is OR connector
func (q *Query) OrWhereOp(key string, operator Operator, value interface{}) *Query {
	return q.AppendCondition(CO_OR, key, operator, value)
}

// OrWhereNil is a function that will add new condition that will check if column nil, connector that is used
// between condition is OR connector
func (q *Query) OrWhereNil(key string) *Query {
	return q.AppendCondition(CO_OR, key, OP_IS_NULL, nil)
}

// OrWhereNotNil is a function that will add new condition that will check if column not nil, connector that is used
// between condition is OR connector
func (q *Query) OrWhereNotNil(key string) *Query {
	return q.AppendCondition(CO_OR, key, OP_IS_NOT_NULL, nil)
}

// GroupWhere is a function that will receive a function that return a group of condition to be grouped
// together, connector that is used between condition is AND connector
func (q *Query) GroupWhere(conditionFunc func(query *Query) *Query) *Query {
	query := QueryTx(q.Tx, q.Config).Use(q.Model)
	param := conditionFunc(query).Parameter
	q.appendGroupConditions(param, CO_AND)
	return q
}

// OrGroupWhere is a function that will receive a function that return a group of condition to be grouped
// together, connector that is used between condition is OR connector
func (q *Query) OrGroupWhere(conditionFunc func(query *Query) *Query) *Query {
	query := QueryTx(q.Tx, q.Config).Use(q.Model)
	param := conditionFunc(query).Parameter
	q.appendGroupConditions(param, CO_OR)
	return q
}

// GroupBy is a function that will add new group by column
func (q *Query) GroupBy(column string) *Query {
	q.appendGroupBy(column)
	return q
}

// OrderBy is a function that will add new order by column with direction
func (q *Query) OrderBy(key string, direction OrderDirection) *Query {
	q.appendOrderBy(q.Model.GetTable(), key, direction)
	return q
}

// WithLimit is a function that will add limit to result row
func (q *Query) WithLimit(limit int) *Query {
	q.setLimit(limit)
	return q
}

// WithOffset is a function that will add offset to result row
func (q *Query) WithOffset(offset int) *Query {
	q.setOffset(offset)
	return q
}

// WithLock is a function that will add lock to current query
func (q *Query) WithLock(lock Lock) *Query {
	q.setLock(lock)
	return q
}

// AppendCondition is a wrapper function for appendCondition
func (q *Query) AppendCondition(connector Connector, key string, operator Operator, value interface{}) *Query {
	q.appendCondition(q.Model.GetTable(), connector, key, operator, value)
	return q
}

// assignModelSlices is a function that will assign Model attribute based on current model used
// to a slices of model of database results
func (q *Query) assignModelSlices(results interface{}, m *Model) interface{} {
	slice := reflect.ValueOf(results).Elem().Interface()
	sVal := reflect.ValueOf(slice)
	for i := 0; i < sVal.Len(); i++ {
		assigned := q.assignModel(sVal.Index(i).Addr().Interface(), m)
		sVal.Index(i).Set(reflect.ValueOf(assigned))
	}
	return sVal.Interface()
}

// assignModel is a function that will assign Model attribute based on current model used
// to a slices of model of database results
func (q *Query) assignModel(result interface{}, m *Model) interface{} {
	value := reflect.ValueOf(result).Interface()
	model := reflect.ValueOf(&m).Elem()
	complete := reflect.ValueOf(value).Elem()
	complete.FieldByName("Model").Set(model)
	return complete.Interface()
}

// first is a function that will get the one result from a query
func (q *Query) first(withTrash TrashStatus) (interface{}, error) {
	defer q.clearParameter()

	err := q.handleNilModel()
	if nil != err {
		return nil, err
	}

	err = q.autoBegin()
	if nil != err {
		return nil, err
	}

	result, err := q.makeTypeOf(q.Model)
	if nil != err {
		q.autoRollback()
		return nil, err
	}

	q.setLimit(1)
	stmt, args, err := q.prepareSelect(AG_NONE, "", withTrash)
	if nil != err {
		q.autoRollback()
		return nil, err
	}

	err = stmt.Get(result, args)
	if nil != err {
		q.autoRollback()
		return nil, err
	}

	q.autoCommit()
	return q.assignModel(result, q.Model.GetModel()), nil
}

// all is a function that will get multiple result from a query
func (q *Query) all(withTrash TrashStatus) (interface{}, error) {
	defer q.clearParameter()

	err := q.handleNilModel()
	if nil != err {
		return nil, err
	}

	err = q.autoBegin()
	if nil != err {
		return nil, err
	}

	results, err := q.makeSliceOf(q.Model)
	if nil != err {
		q.autoRollback()
		return nil, err
	}

	stmt, args, err := q.prepareSelect(AG_NONE, "", withTrash)
	if nil != err {
		q.autoRollback()
		return nil, err
	}

	err = stmt.Select(results, args)
	if nil != err {
		q.autoRollback()
		return nil, err
	}

	q.autoCommit()
	return q.assignModelSlices(results, q.Model.GetModel()), nil
}

// aggregate is a function that will return aggregate value of a column
func (q *Query) aggregate(aggregate Aggregate, column string, withTrash TrashStatus) (*float64, error) {
	defer q.clearParameter()

	var result float64

	err := q.handleNilModel()
	if nil != err {
		return nil, err
	}

	err = q.autoBegin()
	if nil != err {
		return nil, err
	}

	stmt, args, err := q.prepareSelect(aggregate, column, NO_TRASH)
	if nil != err {
		q.autoRollback()
		return nil, err
	}

	err = stmt.Get(&result, args)
	if nil != err {
		q.autoRollback()
		return nil, err
	}

	q.autoCommit()
	return &result, nil
}

// prepareSelect is a function that will return query statement as NamedStmt and parsed payload as a map[string]interface
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

// bindIn is a function that will bind value with named arguments inside in condition
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

// setPKCondition is a function that will append a condition that will check if primary key equals
// current model primary key value, used in Update and Delete method
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

// mergedPayload is a function that will merge model payload with condition values saved in Parameter.Values
func (q *Query) mergedPayload() map[string]interface{} {
	payload := q.Model.Payload()
	for i, v := range q.Parameter.Values {
		payload[i] = v
	}
	return payload
}

// bulkPayload is a function that will merge all payload for bulkinsert into sequential slice of
// map[string]interface{}
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

// autoBegin is a function that will automatically begin a transaction for a query if the Query instance
// is set using sqlx.DB not sqlx.Tx
func (q *Query) autoBegin() error {
	var err error
	if q.AutoCommit && nil != q.Db {
		q.Tx, err = q.Db.Beginx()
		return err
	}
	return nil
}

// autoCommit is a function that will automatically commit a query if the Query instance is set using
// sqlx.DB not sqlx.Tx
func (q *Query) autoCommit() {
	if q.AutoCommit {
		_ = q.Tx.Commit()
	}
}

// autoRollback is a function that will automatically rollback a query if the Query instance is set using
// sqlx.DB not sqlx.Tx
func (q *Query) autoRollback() {
	if q.AutoCommit {
		_ = q.Tx.Rollback()
	}
}

// handleNilModel is a function that will return error if Model attribute is nil, please use Use(v interface{})
// method to set the Model attribute
func (q *Query) handleNilModel() error {
	if nil == q.Model {
		return fmt.Errorf("please use a model before doing query")
	}
	return nil
}

// clearParameter is a function that will clear all condition to prepare query for the next use
func (q *Query) clearParameter() {
	q.Parameter = NewParameter(q.Config)
}

// makeTypeOf is a function to create a new instance of sample Type to make First method immutable
func (q *Query) makeTypeOf(sample interface{}) (interface{}, error) {
	if reflect.TypeOf(sample).Kind() != reflect.Ptr {
		return nil, fmt.Errorf("sample must be a pointer to reference model")
	}
	element := reflect.TypeOf(sample).Elem()
	return reflect.New(element).Interface(), nil
}

// makeSliceOf is a function to create a new instance of sample Type to make All method immutable
func (q *Query) makeSliceOf(sample interface{}) (interface{}, error) {
	if reflect.TypeOf(sample).Kind() != reflect.Ptr {
		return nil, fmt.Errorf("sample must be a pointer to reference model")
	}
	element := reflect.TypeOf(sample).Elem()
	return reflect.New(reflect.SliceOf(element)).Interface(), nil
}
