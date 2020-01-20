package prometheusclient

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/payfazz/go-apt/pkg/fazzrouter"
	"github.com/prometheus/client_golang/prometheus"
)

var httpDurationHistogramOnce sync.Once
var httpDurationHistogram *prometheus.HistogramVec

func httpRequestDurationHistogram() *prometheus.HistogramVec {
	httpDurationHistogramOnce.Do(func() {
		httpDurationHistogram = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_milliseconds",
				Help:    "A histogram of latencies for requests in millisecond.",
				Buckets: []float64{50, 100, 250, 500, 1000, 5000, 10000},
			},
			[]string{"service", "path", "method", "code"},
		)

		prometheus.MustRegister(httpDurationHistogram)
	})

	return httpDurationHistogram
}

// HTTPRequestDurationMiddleware middleware wrapper for ObserveRequestDuration, recommended to be used if you are using `go-apt/pkg/fazzrouter` package, the only thing required: before using this middleware make sure you use `kv.New()` middleware from `github.com/payfazz/go-middleware`
// required params:
// - productName: your product / team name (snake_case)
// - serviceName: your service name (snake_case)
func HTTPRequestDurationMiddleware(productName string, serviceName string) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(writer http.ResponseWriter, req *http.Request) {
			start := time.Now()
			prometheusWriter := wrapResponseWriter(writer)

			next(prometheusWriter, req)

			ObserveRequestDuration(
				productName,
				serviceName,
				fazzrouter.GetPattern(req),
				req.Method,
				prometheusWriter.Code(),
				start,
			)
		}
	}
}

// ObserveRequestDuration observe request duration from startRequestAt until now and store it as histogram, usage example can be seen in HTTPRequestDurationMiddleware method
// required params:
// - productName: your product / team name (snake_case)
// - serviceName: your service name (snake_case)
// - pattern: your route pattern not the requested url, ex: `/v1/users/:id` (correct); `/v1/users/{id}` (correct); `/v1/users/1` (incorrect)
// - method: your http request method (GET, POST, PATCH, etc)
// - code: your http status code (200, 404, 500, etc)
// - startRequestAt: your request start time, recommended to be initiated before `next` method call
func ObserveRequestDuration(productName string, serviceName string, pattern string, method string, code string, startRequestAt time.Time) {
	service := fmt.Sprintf("%s_%s", productName, serviceName)
	labels := prometheus.Labels{
		"service": service,
		"path":    pattern,
		"method":  method,
		"code":    code,
	}
	duration := float64(time.Since(startRequestAt).Milliseconds())
	httpRequestDurationHistogram().With(labels).Observe(duration)
}
