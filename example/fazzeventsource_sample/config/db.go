package config

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"sync"
)

const (
	HOST     = "localhost"
	PORT     = 5432
	USER     = "esexample"
	PASSWORD = "esexample"
	DBNAME   = "esexample"
)

var db *sqlx.DB
var dbOnce sync.Once

// GetDB is a function to get DB instance
func GetDB() *sqlx.DB {
	dbOnce.Do(func() {
		var err error
		conn := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			HOST,
			PORT,
			USER,
			PASSWORD,
			DBNAME,
		)
		db, err = sqlx.Connect("postgres", conn)

		if err != nil {
			panic(err)
		}
	})
	return db
}
