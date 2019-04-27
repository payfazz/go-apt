package fazzdb

import (
	"github.com/jmoiron/sqlx"
)

// Run is a function that used to run the service under tx
func Run(db *sqlx.DB, config Config, fn func(q *Query) error) error {
	tx, err := db.Beginx()
	if nil != err {
		return err
	}

	q := QueryTx(tx, config)
	err = fn(q)
	if nil != err {
		_ = q.Tx.Rollback()
		return err
	}

	_ = q.Tx.Commit()
	return nil
}

// RunDefault basic boiler plate to start the transaction
func RunDefault(db *sqlx.DB, fn func(q *Query) error) error {
	return Run(db, DEFAULT_QUERY_CONFIG, fn)
}
