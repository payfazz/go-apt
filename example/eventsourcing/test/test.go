package test

import (
	"context"
	"github.com/payfazz/go-apt/example/eventsourcing/config"
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

// PrepareTestContenxt is a function to prepare context for testing
func PrepareTestContext() context.Context {
	queryDb := fazzdb.QueryDb(config.GetDB(),
		fazzdb.Config{
			Limit:           20,
			Offset:          0,
			Lock:            fazzdb.LO_NONE,
			DevelopmentMode: true,
		})

	ctx := context.Background()
	ctx = fazzdb.NewQueryContext(ctx, queryDb)

	return ctx
}
