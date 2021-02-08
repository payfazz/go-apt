package prometheusclient

import (
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
)

var pgConnOnce sync.Once
var pgConnectionIdleCount *prometheus.GaugeVec
var pgConnectionUseCount *prometheus.GaugeVec
var pgConnectionWaitCount *prometheus.GaugeVec

// PGConnectionGauge monitor the pooling connection database inside application
func PGConnectionGauge(labels prometheus.Labels, db *sqlx.DB) {
	pgConnOnce.Do(func() {
		pgConnectionIdleCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "pg_connection_idle_count",
			Help: "show the database connection idle",
		}, GetRequiredDBLabels())

		pgConnectionUseCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "pg_connection_use_count",
			Help: "show the database connection use count",
		}, GetRequiredDBLabels())

		pgConnectionWaitCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "pg_connection_wait_count",
			Help: "show the database connection wait count",
		}, GetRequiredDBLabels())

		prometheus.MustRegister(pgConnectionIdleCount, pgConnectionUseCount, pgConnectionWaitCount)
	})

	if !IsValidRequiredDBLabels(labels) {
		return
	}

	pgConnectionIdleCount.With(labels).Set(float64(db.Stats().Idle))
	pgConnectionUseCount.With(labels).Set(float64(db.Stats().InUse))
	pgConnectionWaitCount.With(labels).Set(float64(db.Stats().WaitCount))
}
