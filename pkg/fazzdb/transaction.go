package fazzdb

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// Run is a function that used to run the service under tx
func Run(ctx context.Context, db *sqlx.DB, config Config, fn func() error) error {
	tx, err := db.Beginx()
	if nil != err {
		return err
	}

	q := QueryTx(tx, config)
	ctx = NewTransactionContext(ctx, q)

	err = fn()
	if nil != err {
		_ = q.Tx.Rollback()
		return err
	}

	_ = q.Tx.Commit()
	return nil
}

// RunDefault basic boiler plate to start the transaction
func RunDefault(ctx context.Context, db *sqlx.DB, fn func() error) error {
	return Run(ctx, db, DEFAULT_QUERY_CONFIG, fn)
}
