package test

import (
	"context"
	"github.com/payfazz/go-apt/config"
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

// PrepareTestContenxt is a function to prepare context for testing
func PrepareTestContext() context.Context {
	queryDb := fazzdb.QueryDb(config.GetDB(), config.Parameter)
	ctx := context.Background()
	ctx = fazzdb.NewQueryContext(ctx, queryDb)
	return ctx
}
