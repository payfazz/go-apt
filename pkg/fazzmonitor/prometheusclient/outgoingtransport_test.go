package prometheusclient

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/require"
)

func TestOutgoingHTTPTransportWithMetrics(t *testing.T) {
	scenarios := map[string]struct {
		method                  string
		url                     string
		body                    io.Reader
		transport               http.RoundTripper
		expectedInflightMetrics string
	}{
		"default transport": {
			http.MethodGet,
			"https://api.ipify.org/?format=text",
			nil,
			nil,
			"",
		},
		"use outgoing http metrics but disabled": {
			http.MethodGet,
			"https://api.ipify.org/?format=text",
			nil,
			OutgoingHTTPTransportWithMetrics(false, nil),
			"",
		},
		"use outgoing http metrics": {
			http.MethodGet,
			"https://api.ipify.org/?format=text",
			nil,
			OutgoingHTTPTransportWithMetrics(true, &http.Transport{}),
			`
			# HELP outgoing_http_request_inflight_count outgoing requests that have been submitted but have not been completed
			# TYPE outgoing_http_request_inflight_count gauge
			outgoing_http_request_inflight_count{host="api.ipify.org",method="GET",path="/",protocol="https"} 0
			`,
		},
	}

	for name, s := range scenarios {
		t.Run(name, func(t *testing.T) {
			c := http.Client{
				Transport: s.transport,
			}

			req, err := http.NewRequest(s.method, s.url, s.body)
			if err != nil {
				t.Fatal(err)
			}

			_, err = c.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			require.Nil(t, testutil.CollectAndCompare(outgoingHTTPInflightCount, strings.NewReader(s.expectedInflightMetrics)))
		})
	}
}
