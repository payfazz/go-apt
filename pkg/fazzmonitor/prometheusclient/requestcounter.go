package prometheusclient

import (
	"net/http"
	"strconv"
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
			[]string{"path", "method", "code"},
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
// - code: your http status code (200, 404, 500, etc)
func IncrementRequestCounter(pattern string, method string, code string) {
	labels := prometheus.Labels{
		"path":   pattern,
		"method": method,
		"code":   code,
	}
	httpRequestCounter().With(labels).Inc()
}

type interceptingWriter struct {
	http.ResponseWriter
	code int
}

func (w *interceptingWriter) WriteHeader(code int) {
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}

// InstrumentHandlerCounter is a decorator that wraps the provided http.Handler
// to use IncrementRequestCounter after wrapped handler invoked.
//
// If wrapped handler does not set a status code, a status code 200 is assumed.
func InstrumentHandlerCounter(pattern string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		iw := &interceptingWriter{w, http.StatusOK}
		h.ServeHTTP(iw, r)
		IncrementRequestCounter(pattern, r.Method, strconv.Itoa(iw.code))
	})
}
