package fazzdb

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Run is a function that used to run the service under tx
func Run(ctx context.Context, db *sqlx.DB, config Config, fn func(ctx context.Context) error) error {
	tx, err := db.BeginTxx(ctx, config.Opts)
	if nil != err {
		return err
	}

	q := QueryTx(tx, config)
	ctx = NewTransactionContext(ctx, q)

	err = fn(ctx)
	if nil != err {
		_ = q.Tx.Rollback()
		return err
	}

	_ = q.Tx.Commit()
	return nil
}

// RunDefault basic boiler plate to start the transaction
func RunDefault(ctx context.Context, db *sqlx.DB, fn func(ctx context.Context) error) error {
	return Run(ctx, db, DEFAULT_QUERY_CONFIG, fn)
}

// GetTxOptions is a function that will create default TxOptions if given opts is nil
func GetTxOptions(opts *sql.TxOptions) *sql.TxOptions {
	if nil == opts {
		return &sql.TxOptions{
			Isolation: sql.LevelDefault,
			ReadOnly:  false,
		}
	}

	return opts
}
