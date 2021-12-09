package fazzrepository

import (
	"context"
	"reflect"

	"github.com/payfazz/go-apt/pkg/fazzdb"
)

// RepositoryInterface interface for repository struct
type RepositoryInterface interface {
	GetQuery(ctx context.Context) (*fazzdb.Query, error)

	Count(ctx context.Context, conditions []fazzdb.SliceCondition) (*float64, error)
	Sum(ctx context.Context, column string, conditions []fazzdb.SliceCondition) (*float64, error)

	RawSelect(ctx context.Context, sample interface{}, query string, payload ...interface{}) (interface{}, error)
	RawExec(ctx context.Context, query string, payload ...interface{}) (bool, error)

	FindAll(ctx context.Context, conditions []fazzdb.SliceCondition, orders []fazzdb.Order, limit int, offset int) (interface{}, error)
	FindOne(ctx context.Context, conditions []fazzdb.SliceCondition, orders []fazzdb.Order) (interface{}, error)
	Find(ctx context.Context, id interface{}) (interface{}, error)

	Create(ctx context.Context, m fazzdb.ModelInterface) (interface{}, error)
	Update(ctx context.Context, m fazzdb.ModelInterface) (bool, error)
	Delete(ctx context.Context, m fazzdb.ModelInterface) (bool, error)
}

// Repository base struct for all repository
type Repository struct {
	model              fazzdb.ModelInterface
	withoutTransaction bool
}

// GetQuery get query instance from context
func (r *Repository) GetQuery(ctx context.Context) (*fazzdb.Query, error) {
	if !r.withoutTransaction {
		return fazzdb.GetTransactionOrQueryContext(ctx)
	}
	return fazzdb.GetQueryContext(ctx)
}

// Count find count of rows using given conditions
func (r *Repository) Count(ctx context.Context, conditions []fazzdb.SliceCondition) (*float64, error) {
	q, err := r.GetQuery(ctx)
	if nil != err {
		return nil, err
	}

	return q.Use(r.model).
		WhereMany(conditions...).
		CountCtx(ctx)
}

// Sum find sum of column using given conditions
func (r *Repository) Sum(ctx context.Context, column string, conditions []fazzdb.SliceCondition) (*float64, error) {
	q, err := r.GetQuery(ctx)
	if nil != err {
		return nil, err
	}

	return q.Use(r.model).
		WhereMany(conditions...).
		SumCtx(ctx, column)
}

// RawSelect find data by raw query and payload if there is arguments
func (r *Repository) RawSelect(ctx context.Context, sample interface{}, query string, payload ...interface{}) (interface{}, error) {
	q, err := r.GetQuery(ctx)
	if nil != err {
		return nil, err
	}

	return q.RawAllCtx(ctx, sample, query, payload...)
}

// RawExec execute query other than SELECT and return success status of the query
func (r *Repository) RawExec(ctx context.Context, query string, payload ...interface{}) (bool, error) {
	q, err := r.GetQuery(ctx)
	if nil != err {
		return false, err
	}

	return q.RawExecCtx(ctx, query, payload...)
}

// FindAll find data by given conditions, order, limit and offset
func (r *Repository) FindAll(
	ctx context.Context,
	conditions []fazzdb.SliceCondition,
	orders []fazzdb.Order,
	limit int,
	offset int,
) (interface{}, error) {
	q, err := r.GetQuery(ctx)
	if nil != err {
		return nil, err
	}

	current := q.Use(r.model).
		WhereMany(conditions...).
		OrderByMany(orders...)

	if limit > 0 {
		current.WithLimit(limit)
	} else if limit == NO_LIMIT {
		current.WithLimit(0)
	}

	if offset > 0 {
		current.WithOffset(offset)
	}

	return current.AllCtx(ctx)
}

// FindOne find one data by given conditions and orders
func (r *Repository) FindOne(
	ctx context.Context,
	conditions []fazzdb.SliceCondition,
	orders []fazzdb.Order,
) (interface{}, error) {
	q, err := r.GetQuery(ctx)
	if nil != err {
		return nil, err
	}

	current := q.Use(r.model).
		WhereMany(conditions...).
		OrderByMany(orders...).
		WithLimit(1)

	rows, err := current.AllCtx(ctx)
	if nil != err {
		return nil, err
	}

	val := reflect.ValueOf(rows)
	if val.Len() == 0 {
		return nil, NewEmptyResultError()
	}

	return val.Index(0).Interface(), nil
}

// Find find data by given id
func (r *Repository) Find(ctx context.Context, id interface{}) (interface{}, error) {
	q, err := r.GetQuery(ctx)
	if nil != err {
		return nil, err
	}

	rows, err := q.Use(r.model).
		Where("id", id).
		WithLimit(1).
		AllCtx(ctx)

	if nil != err {
		return nil, err
	}

	val := reflect.ValueOf(rows)
	if val.Len() == 0 {
		return nil, NewEmptyResultError()
	}

	return val.Index(0).Interface(), nil
}

// Create insert data by given model
func (r *Repository) Create(ctx context.Context, m fazzdb.ModelInterface) (interface{}, error) {
	q, err := r.GetQuery(ctx)
	if nil != err {
		return nil, err
	}

	result, err := q.Use(m).
		InsertCtx(ctx, false)

	if nil != err {
		return nil, err
	}

	return result, nil
}

// Update update data by given model
func (r *Repository) Update(ctx context.Context, m fazzdb.ModelInterface) (bool, error) {
	q, err := r.GetQuery(ctx)
	if nil != err {
		return false, err
	}

	_, err = q.Use(m).
		UpdateCtx(ctx)

	if nil != err {
		return false, err
	}

	return true, nil
}

// Delete delete data by given model
func (r *Repository) Delete(ctx context.Context, m fazzdb.ModelInterface) (bool, error) {
	q, err := r.GetQuery(ctx)
	if nil != err {
		return false, err
	}

	_, err = q.Use(m).
		DeleteCtx(ctx)

	if nil != err {
		return false, err
	}

	return true, nil
}

// NewRepository construct base repository
func NewRepository(m fazzdb.ModelInterface) RepositoryInterface {
	return &Repository{
		model: m,
	}
}

// NewRepositoryWithoutTransaction construct base repository without using DB Transaction
func NewRepositoryWithoutTransaction(m fazzdb.ModelInterface) RepositoryInterface {
	return &Repository{
		model:              m,
		withoutTransaction: true,
	}
}
