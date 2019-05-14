package fazzdb

import (
	"context"
)

type txKeyType struct{}
type qdbKeyType struct{}

var txKey txKeyType
var qdbKey qdbKeyType

// Context is a struct to handle query object inside context
type Context struct {
	Query *Query
}

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
func GetTransactionContext(ctx context.Context) *Query {
	return ctx.Value(txKey).(*Query)
}

// GetQueryContext is a function to get db query object from context.
// Must be used after NewQueryContext
func GetQueryContext(ctx context.Context) *Query {
	return ctx.Value(qdbKey).(*Query)
}
