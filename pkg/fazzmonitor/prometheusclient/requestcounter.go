package prometheusclient

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

func httpRequestCounter() *prometheus.CounterVec {
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "A counter for requests to the wrapped handler.",
		},
		[]string{"service", "path", "method", "code"},
	)

	prometheus.MustRegister(counter)

	return counter
}

func RequestCounter(serviceName string, pattern RoutePattern) func(next http.HandlerFunc) http.HandlerFunc {
	counter := httpRequestCounter()

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(writer http.ResponseWriter, req *http.Request) {
			prometheusWriter := wrapResponseWriter(writer)
			next(prometheusWriter, req)
			counter.With(labels(serviceName, prometheusWriter, req, pattern)).Inc()
		}
	}
}
