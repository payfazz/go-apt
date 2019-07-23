package esfazz

import (
	"context"
	"errors"
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

// GetQuery get query is a function that used to get the fazzdb.Query
func getContext(ctx context.Context) (*fazzdb.Query, error) {
	q, _ := fazzdb.GetTransactionContext(ctx)
	if q != nil {
		return q, nil
	}

	qb, _ := fazzdb.GetQueryContext(ctx)
	if qb != nil {
		return qb, nil
	}

	return nil, errors.New("can't find both queryDb and queryTx instance on context")
}
