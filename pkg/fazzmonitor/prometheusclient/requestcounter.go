package prometheusclient

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var httpCounterOnce sync.Once
var httpCounter *prometheus.CounterVec

func httpRequestCounter() *prometheus.CounterVec {
	httpCounterOnce.Do(func() {
		httpCounter = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "A counter for requests to the wrapped handler.",
			},
			[]string{"path", "method", "status"},
		)

		prometheus.MustRegister(httpCounter)
	})

	return httpCounter
}

// IncrementRequestCounter increment http request count and store it as total requests, usage example can be seen in HTTPRequestCounterMiddleware method
// required params:
// - serviceName: your service name (snake_case)
// - pattern: your route pattern not the requested url, ex: `/v1/users/:id` (correct); `/v1/users/{id}` (correct); `/v1/users/1` (incorrect)
// - method: your http request method (GET, POST, PATCH, etc)
// - status: your http status code (200, 404, 500, etc)
func IncrementRequestCounter(pattern string, method string, status string) {
	labels := prometheus.Labels{
		"path":   pattern,
		"method": method,
		"status": status,
	}
	httpRequestCounter().With(labels).Inc()
}
