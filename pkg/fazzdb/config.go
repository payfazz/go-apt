package fazzdb

import (
	"database/sql"
)

// Config is a struct that will be used to set default value for some parameter attribute
type Config struct {
	Limit           int
	Offset          int
	Lock            Lock
	DevelopmentMode bool
	Opts            *sql.TxOptions
}
