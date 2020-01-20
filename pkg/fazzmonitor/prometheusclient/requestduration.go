package prometheusclient

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func httpRequestDurationSummary() *prometheus.SummaryVec {
	summary := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_request_duration_summary",
			Help: "Request latency distributions in milliseconds.",
		},
		[]string{"service", "path", "method", "code"},
	)

	prometheus.MustRegister(summary)

	return summary
}

func httpRequestDurationHistogram() *prometheus.HistogramVec {
	histogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_histogram",
			Help:    "A histogram of latencies for requests in millisecond.",
			Buckets: []float64{10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000},
		},
		[]string{"service", "path", "method", "code"},
	)

	prometheus.MustRegister(histogram)

	return histogram
}

func RequestDuration(serviceName string, pattern RoutePattern) func(next http.HandlerFunc) http.HandlerFunc {
	summary := httpRequestDurationSummary()
	histogram := httpRequestDurationHistogram()

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(writer http.ResponseWriter, req *http.Request) {
			start := time.Now()
			prometheusWriter := wrapResponseWriter(writer)

			next(prometheusWriter, req)

			duration := float64(time.Since(start).Milliseconds())

			summary.With(labels(serviceName, prometheusWriter, req, pattern)).Observe(duration)
			histogram.With(labels(serviceName, prometheusWriter, req, pattern)).Observe(duration)
		}
	}
}
