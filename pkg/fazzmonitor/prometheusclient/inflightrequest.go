package prometheusclient

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var httpGaugeOnce sync.Once
var httpGauge *prometheus.GaugeVec

func httpInflightRequest() *prometheus.GaugeVec {
	httpGaugeOnce.Do(func() {
		httpGauge = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "http_request_inflight_count",
				Help: "Total requests that have been submitted but have not been completed.",
			},
			[]string{"path", "method"},
		)

		prometheus.MustRegister(httpGauge)
	})

	return httpGauge
}

// IncrementInflightRequest increment http inflight request count and store it as total inflight requests
// required params:
// - pattern: your route pattern not the requested url, ex: `/v1/users/:id` (correct); `/v1/users/{id}` (correct); `/v1/users/1` (incorrect)
// - method: your http request method (GET, POST, PATCH, etc)
func IncrementInflightRequest(pattern string, method string) {
	labels := prometheus.Labels{
		"path":   pattern,
		"method": method,
	}
	httpInflightRequest().With(labels).Inc()
}

// DecrementInflightRequest decrement http inflight request count and store it as total inflight requests
// required params:
// - pattern: your route pattern not the requested url, ex: `/v1/users/:id` (correct); `/v1/users/{id}` (correct); `/v1/users/1` (incorrect)
// - method: your http request method (GET, POST, PATCH, etc)
func DecrementInflightRequest(pattern string, method string) {
	labels := prometheus.Labels{
		"path":   pattern,
		"method": method,
	}
	httpInflightRequest().With(labels).Dec()
}
