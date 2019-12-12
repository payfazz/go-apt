package prometheusclient

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func RequestDuration() func(next http.HandlerFunc) http.HandlerFunc {
	collector := registerOnce(
		prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name: "http_request_duration_seconds",
				Help: "Request latency distributions.",
			},
			[]string{"date", "metrics", "method", "path"},
		),
	).(*prometheus.SummaryVec)

	pathHistogramCollector := registerOnce(
		prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds_histogram",
				Help:    "A histogram of latencies for requests in millisecond.",
				Buckets: []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000},
			},
			[]string{"date", "method", "path"},
		),
	).(*prometheus.HistogramVec)

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next(w, r)

			duration := float64(time.Since(start).Milliseconds())

			collector.WithLabelValues(
				dateMinute(),
				"duration",
				r.Method,
				path(r),
			).Observe(duration)

			pathHistogramCollector.WithLabelValues(
				dateMinute(),
				r.Method,
				path(r),
			).Observe(duration)
		}
	}
}
