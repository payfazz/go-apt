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
			[]string{"metrics", "method", "path"},
		),
	).(*prometheus.SummaryVec)

	histogramCollector := registerOnce(
		prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds_histogram",
				Help:    "A histogram of latencies for requests in millisecond.",
				Buckets: []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000},
			},
		),
	).(prometheus.Histogram)

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next(w, r)

			duration := float64(time.Since(start).Milliseconds())

			collector.WithLabelValues(
				"duration",
				r.Method,
				path(r),
			).Observe(duration)

			histogramCollector.Observe(duration)
		}
	}
}
