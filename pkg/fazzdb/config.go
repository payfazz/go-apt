package fazzdb

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

// Config is a struct that will be used to set default value for some parameter attribute
type Config struct {
	Limit           int
	Offset          int
	Lock            Lock
	DevelopmentMode bool
	PrometheusMode  bool
	Labels          prometheus.Labels
	Opts            *sql.TxOptions
}
