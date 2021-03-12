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

// DBQueryMetrics monitor the query execution includes latency, error and inflight
func DBQueryMetrics(labels prometheus.Labels, query string, prometheusMode bool, fn func() error) error {
	queryOnce.Do(func() {
		pgQueryInflightCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "pg_query_inflight_count",
			Help: "requests that have been submitted but have not been completed.",
		}, append([]string{"query"}, GetRequiredDBLabels()...))
		pgQueryDurationSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "pg_query_duration_seconds",
			Help:    "latency of query execution",
			Buckets: prometheus.DefBuckets,
		}, append([]string{"query"}, GetRequiredDBLabels()...))
		pgQueryErrorCount = prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "pg_query_error_count",
			Help: "error count when querying to db",
		}, append([]string{"query"}, GetRequiredDBLabels()...))

		prometheus.MustRegister(pgQueryInflightCount, pgQueryDurationSeconds, pgQueryErrorCount)
	})

	if !prometheusMode || !IsValidRequiredDBLabels(labels) {
		return fn()
	}

	labels["query"] = query

	start := time.Now()
	pgQueryInflightCount.With(labels).Inc()
	defer pgQueryInflightCount.With(labels).Dec()

	err := fn()
	if err != nil {
		pgQueryErrorCount.With(labels).Inc()
	}

	duration := time.Since(start)
	pgQueryDurationSeconds.With(labels).Observe(duration.Seconds())

	return err
}
