package prometheusclient

import (
	"net/http"

	"github.com/payfazz/go-apt/pkg/fazzcommon/formatter"
	"github.com/payfazz/go-apt/pkg/fazzmonitor"
	"github.com/prometheus/client_golang/prometheus"
)

func StatusCounter() func(next http.HandlerFunc) http.HandlerFunc {
	collector := registerOnce(
		prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_header_status",
				Help: "Total number of scrapes by HTTP status code.",
			},
			[]string{"date", "method", "path", "code"},
		),
	).(*prometheus.CounterVec)

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			var wr *fazzmonitor.Writer
			var ok bool

			if wr, ok = w.(*fazzmonitor.Writer); !ok {
				wr = &fazzmonitor.Writer{ResponseWriter: w}
			}

			next(wr, r)

			collector.WithLabelValues(
				dateMinute(),
				path(r),
				r.Method,
				formatter.IntegerToString(wr.StatusCode),
			).Inc()
		}
	}
}
