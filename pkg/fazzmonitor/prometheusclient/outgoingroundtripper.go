package prometheusclient

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var outgoingHTTPRoundTripperOnce sync.Once
var outgoingHTTPRoundTripperInflightCount *prometheus.GaugeVec
var outgoingHTTPRoundTripperDurationSeconds *prometheus.HistogramVec
var outgoingHTTPRoundTripperRequestsTotal *prometheus.CounterVec

type metricsRoundTripper struct {
	base http.RoundTripper
}

// RoundTrip satisfy RoundTripper interface from http package
func (m *metricsRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	inflightLabels := prometheus.Labels{
		"host":     req.URL.Host,
		"path":     req.URL.Path,
		"method":   strings.ToUpper(req.Method),
		"protocol": req.URL.Scheme,
	}
	outgoingHTTPRoundTripperInflightCount.With(inflightLabels).Inc()
	defer outgoingHTTPRoundTripperInflightCount.With(inflightLabels).Dec()

	start := time.Now()
	res, err := m.base.RoundTrip(req)
	durationSeconds := time.Since(start).Seconds()

	durationLabels := prometheus.Labels{
		"host":     req.URL.Host,
		"path":     req.URL.Path,
		"method":   strings.ToUpper(req.Method),
		"protocol": req.URL.Scheme,
		"code":     "",
	}
	if res != nil {
		durationLabels["code"] = strconv.Itoa(res.StatusCode)
	}

	outgoingHTTPRoundTripperDurationSeconds.With(durationLabels).Observe(durationSeconds)

	requestCountLabels := prometheus.Labels{
		"host":     req.URL.Host,
		"path":     req.URL.Path,
		"method":   strings.ToUpper(req.Method),
		"protocol": req.URL.Scheme,
		"code":     "",
	}
	if res != nil {
		requestCountLabels["code"] = strconv.Itoa(res.StatusCode)
	}
	outgoingHTTPRoundTripperRequestsTotal.With(requestCountLabels).Inc()

	return res, err
}

// OutgoingHTTPRoundTripperWithMetrics wrap the http RoundTripper to provide outgoing http metrics
func OutgoingHTTPRoundTripperWithMetrics(enable bool, base http.RoundTripper) http.RoundTripper {
	outgoingHTTPRoundTripperOnce.Do(func() {
		outgoingHTTPRoundTripperInflightCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "outgoing_http_request_inflight_count",
			Help: "outgoing requests that have been submitted but have not been completed",
		}, []string{"host", "path", "method", "protocol"})
		outgoingHTTPRoundTripperDurationSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "outgoing_http_request_duration_seconds",
			Help:    "latency of the outgoing requests.",
			Buckets: prometheus.DefBuckets,
		}, []string{"host", "path", "method", "protocol", "code"})
		outgoingHTTPRoundTripperRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "outgoing_http_requests_total",
			Help: "Total number of HTTP requests completed on the server, regardless of success or failure",
		}, []string{"host", "path", "method", "protocol", "code"})

		prometheus.MustRegister(outgoingHTTPRoundTripperInflightCount, outgoingHTTPRoundTripperDurationSeconds, outgoingHTTPRoundTripperRequestsTotal)
	})

	if base == nil {
		base = &http.Transport{}
	}

	if !enable {
		return base
	}

	m := metricsRoundTripper{base}
	return &m
}
