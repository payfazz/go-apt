package prometheusclient

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

func RequestCounter() func(next http.HandlerFunc) http.HandlerFunc {
	collector := registerOnce(
		prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "A counter for requests to the wrapped handler.",
			},
			[]string{"date", "method", "path"},
		),
	).(*prometheus.CounterVec)

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			next(w, r)

			collector.WithLabelValues(
				dateMinute(),
				path(r),
				r.Method,
			).Inc()
		}
	}
}
