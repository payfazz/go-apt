package prometheusclient

import (
	"fmt"
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
			[]string{"service", "path", "method", "code"},
		)

		prometheus.MustRegister(httpCounter)
	})

	return httpCounter
}

// IncrementRequestCounter increment http request count and store it as total requests, usage example can be seen in HTTPRequestCounterMiddleware method
// required params:
// - productName: your product / team name (snake_case)
// - serviceName: your service name (snake_case)
// - pattern: your route pattern not the requested url, ex: `/v1/users/:id` (correct); `/v1/users/{id}` (correct); `/v1/users/1` (incorrect)
// - method: your http request method (GET, POST, PATCH, etc)
// - code: your http status code (200, 404, 500, etc)
func IncrementRequestCounter(productName string, serviceName string, pattern string, method string, code string) {
	service := fmt.Sprintf("%s_%s", productName, serviceName)
	labels := prometheus.Labels{
		"service": service,
		"path":    pattern,
		"method":  method,
		"code":    code,
	}
	httpRequestCounter().With(labels).Inc()
}
