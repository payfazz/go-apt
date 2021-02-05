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

func PGConnectionGauge(labels prometheus.Labels, db *sqlx.DB) {
	if !IsValidRequiredDBLabels(labels) {
		return
	}

	pgConnOnce.Do(func() {
		pgConnectionIdleCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "pg_connection_idle_count",
			Help: "show the count of database connection idle",
		}, GetRequiredDBLabels())

		pgConnectionUseCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "pg_connection_use_count",
			Help: "show the count of database connection used count",
		}, GetRequiredDBLabels())

		pgConnectionWaitCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "pg_connection_wait_count",
			Help: "show the count of database waiting for a new connection count",
		}, GetRequiredDBLabels())

		prometheus.MustRegister(pgConnectionIdleCount, pgConnectionUseCount, pgConnectionWaitCount)
	})

	pgConnectionIdleCount.With(labels).Set(float64(db.Stats().Idle))
	pgConnectionUseCount.With(labels).Set(float64(db.Stats().InUse))
	pgConnectionWaitCount.With(labels).Set(float64(db.Stats().WaitCount))
}
