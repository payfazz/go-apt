package fazzdb

import (
	"context"
	"errors"
)

type txKeyType struct{}
type qdbKeyType struct{}

var txKey txKeyType
var qdbKey qdbKeyType

// NewTransactionContext is a function to append transaction query object into context
func NewTransactionContext(ctx context.Context, queryTx *Query) context.Context {
	return context.WithValue(ctx, txKey, queryTx)
}

// NewQueryContext is a function to append db query object into context
func NewQueryContext(ctx context.Context, queryDb *Query) context.Context {
	return context.WithValue(ctx, qdbKey, queryDb)
}

// GetTransactionContext is a function to get transaction query object from context
// Must be used after NewTransactionContext
func GetTransactionContext(ctx context.Context) (*Query, error) {
	query := ctx.Value(txKey)
	if nil == query {
		return nil, errors.New("transaction instance not found in context, make sure to call NewTransactionContext before calling GetTransactionContext")
	}
	return query.(*Query), nil
}

// GetQueryContext is a function to get db query object from context.
// Must be used after NewQueryContext
func GetQueryContext(ctx context.Context) (*Query, error) {
	query := ctx.Value(qdbKey)
	if nil == query {
		return nil, errors.New("query db instance not found in context, make sure to call NewQueryContext before calling GetQueryContext")
	}
	return query.(*Query), nil
}
