package middleware

import (
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/payfazz/go-apt/pkg/fazzcommon/response"
	"github.com/payfazz/go-apt/pkg/fazzmonitor/prometheusclient"
	"github.com/payfazz/go-apt/pkg/fazzrouter"
)

// HTTPRequestCounterMiddleware middleware wrapper for IncrementRequestCounter, recommended to be used if you are using `go-apt/pkg/fazzrouter` package, the only thing required: before using this middleware make sure you use `kv.New()` middleware from `github.com/payfazz/go-middleware`
func HTTPRequestCounterMiddleware() func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(writer http.ResponseWriter, req *http.Request) {
			prometheusWriter := response.WrapWriter(writer)

			next(prometheusWriter, req)

			prometheusclient.IncrementRequestCounter(
				fazzrouter.GetPattern(req),
				req.Method,
				prometheusWriter.Code(),
			)
		}
	}
}

// HTTPRequestDurationMiddleware middleware wrapper for ObserveRequestDuration, recommended to be used if you are using `go-apt/pkg/fazzrouter` package, the only thing required: before using this middleware make sure you use `kv.New()` middleware from `github.com/payfazz/go-middleware`
func HTTPRequestDurationMiddleware() func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(writer http.ResponseWriter, req *http.Request) {
			start := time.Now()
			prometheusWriter := response.WrapWriter(writer)

			next(prometheusWriter, req)

			prometheusclient.ObserveRequestDuration(
				fazzrouter.GetPattern(req),
				req.Method,
				prometheusWriter.Code(),
				start,
			)
		}
	}
}

// PGConnectionMiddleware middleware wrapper for PGConnectionGauge, recommended to be used if you are using `go-apt/pkg/fazzrouter` package, the only thing required: before using this middleware make sure you use `kv.New()` middleware from `github.com/payfazz/go-middleware`
func PGConnectionMiddleware(labels prometheus.Labels, db *sqlx.DB) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(writer http.ResponseWriter, req *http.Request) {
			prometheusclient.PGConnectionGauge(labels, db)

			next(writer, req)
		}
	}
}
