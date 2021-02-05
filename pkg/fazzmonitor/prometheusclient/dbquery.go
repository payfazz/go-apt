package prometheusclient

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var queryOnce sync.Once
var pgQueryInflightCount *prometheus.GaugeVec
var pgQueryErrorCount *prometheus.CounterVec
var pgQueryDurationSeconds *prometheus.HistogramVec

func DBQueryMetrics(labels prometheus.Labels, query string, prometheusMode bool, fn func() error) error {
	if !prometheusMode {
		return fn()
	}

	queryOnce.Do(func() {
		pgQueryInflightCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "pg_query_inflight_count",
			Help: "requests that have been submitted but have not been completed.",
		}, []string{"host", "port", "name", "user", "query"})
		pgQueryDurationSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "pg_query_duration_seconds",
			Help:    "latency of query execution",
			Buckets: []float64{0.1, 0.3, 1, 30, 60},
		}, []string{"host", "port", "name", "user", "query"})
		pgQueryErrorCount = prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "pg_query_error_count",
			Help: "error count when querying to db",
		}, []string{"host", "port", "name", "user", "query"})

		prometheus.MustRegister(pgQueryInflightCount, pgQueryDurationSeconds, pgQueryErrorCount)
	})

	labels["query"] = query

	start := time.Now()
	pgQueryInflightCount.With(labels).Inc()

	err := fn()
	if err != nil {
		pgQueryErrorCount.With(labels).Inc()
	}

	duration := time.Since(start)
	pgQueryInflightCount.With(labels).Dec()
	pgQueryDurationSeconds.With(labels).Observe(duration.Seconds())

	return err
}
