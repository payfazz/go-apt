package prometheusclient

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/payfazz/go-apt/pkg/fazzrouter"
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

// HTTPRequestCounterMiddleware middleware wrapper for IncrementRequestCounter, recommended to be used if you are using `go-apt/pkg/fazzrouter` package, the only thing required: before using this middleware make sure you use `kv.New()` middleware from `github.com/payfazz/go-middleware`
// required params:
// - productName: your product / team name (snake_case)
// - serviceName: your service name (snake_case)
func HTTPRequestCounterMiddleware(productName string, serviceName string) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(writer http.ResponseWriter, req *http.Request) {
			prometheusWriter := wrapResponseWriter(writer)

			next(prometheusWriter, req)

			IncrementRequestCounter(
				productName,
				serviceName,
				fazzrouter.GetPattern(req),
				req.Method,
				prometheusWriter.Code(),
			)
		}
	}
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
