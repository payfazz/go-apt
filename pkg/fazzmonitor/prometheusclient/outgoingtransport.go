package prometheusclient

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var outgoingHTTPTransportOnce sync.Once
var outgoingHTTPInflightCount *prometheus.GaugeVec
var outgoingHTTPDurationSeconds *prometheus.HistogramVec

type metricsRoundTriper struct {
	transport *http.Transport
}

// RoundTrip satisfy RoundTripper interface from http package
func (m *metricsRoundTriper) RoundTrip(req *http.Request) (*http.Response, error) {
	inflightLabels := prometheus.Labels{
		"host":     req.URL.Host,
		"path":     req.URL.Path,
		"method":   strings.ToUpper(req.Method),
		"protocol": req.URL.Scheme,
	}
	outgoingHTTPInflightCount.With(inflightLabels).Inc()
	defer outgoingHTTPInflightCount.With(inflightLabels).Dec()

	start := time.Now()
	res, err := m.transport.RoundTrip(req)
	durationSeconds := time.Since(start).Seconds()

	durationLabels := prometheus.Labels{
		"host":     req.URL.Host,
		"path":     req.URL.Path,
		"method":   strings.ToUpper(req.Method),
		"protocol": req.URL.Scheme,
	}

	if res != nil {
		durationLabels["code"] = strconv.Itoa(res.StatusCode)
	}

	outgoingHTTPDurationSeconds.With(durationLabels).Observe(durationSeconds)
	return res, err
}

// OutgoingHTTPTransportWithMetrics wrap the http transport which implement RoundTripper interface to provide outgoing http metrics
func OutgoingHTTPTransportWithMetrics(enable bool, transport *http.Transport) http.RoundTripper {
	outgoingHTTPTransportOnce.Do(func() {
		outgoingHTTPInflightCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "outgoing_http_request_inflight_count",
			Help: "outgoing requests that have been submitted but have not been completed",
		}, []string{"host", "path", "method", "protocol"})
		outgoingHTTPDurationSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "outgoing_http_request_duration_seconds",
			Help:    "latency of the outgoing requests.",
			Buckets: prometheus.DefBuckets,
		}, []string{"host", "path", "method", "protocol", "code"})

		prometheus.MustRegister(outgoingHTTPInflightCount, outgoingHTTPDurationSeconds)
	})

	if transport == nil {
		transport = &http.Transport{}
	}

	if !enable {
		return transport
	}

	m := metricsRoundTriper{transport}
	return &m
}
