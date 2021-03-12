package prometheusclient

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var httpDurationHistogramOnce sync.Once
var httpDurationHistogram *prometheus.HistogramVec

func httpRequestDurationHistogram() *prometheus.HistogramVec {
	httpDurationHistogramOnce.Do(func() {
		httpDurationHistogram = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "A histogram of latencies for requests in second.",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"path", "method", "status"},
		)

		prometheus.MustRegister(httpDurationHistogram)
	})

	return httpDurationHistogram
}

// ObserveRequestDuration observe request duration from startRequestAt until now and store it as histogram, usage example can be seen in HTTPRequestDurationMiddleware method
// required params:
// - serviceName: your service name (snake_case)
// - pattern: your route pattern not the requested url, ex: `/v1/users/:id` (correct); `/v1/users/{id}` (correct); `/v1/users/1` (incorrect)
// - method: your http request method (GET, POST, PATCH, etc)
// - status: your http status code (200, 404, 500, etc)
// - startRequestAt: your request start time, recommended to be initiated before `next` method call
func ObserveRequestDuration(pattern string, method string, status string, startRequestAt time.Time) {
	labels := prometheus.Labels{
		"path":   pattern,
		"method": method,
		"status": status,
	}
	duration := float64(time.Since(startRequestAt).Seconds())
	httpRequestDurationHistogram().With(labels).Observe(duration)
}
