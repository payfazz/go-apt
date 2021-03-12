package prometheusclient

import (
	"errors"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/require"
)

func TestDBQueryMetrics(t *testing.T) {
	scenarios := map[string]struct {
		err                         bool
		validValidation             bool
		query                       string
		prometheusMode              bool
		labels                      prometheus.Labels
		beforeInflightCountExpected string
		afterInflightCountExpected  string
		durationSecondsExpected     string
		errorCountExpected          string
	}{
		"disable prometheus": {
			prometheusMode: false,
		},
		"invalid validation but prometheus enable": {
			validValidation: false,
			prometheusMode:  true,
		},
		"valid validation but error from query to db": {
			err:             true,
			validValidation: true,
			query:           "SELECT 1",
			prometheusMode:  true,
			labels:          prometheus.Labels{"host": "127.0.0.1", "port": "5432", "user": "postgres", "name": "postgres"},
			beforeInflightCountExpected: `
			# HELP pg_query_inflight_count requests that have been submitted but have not been completed.
			# TYPE pg_query_inflight_count gauge
			pg_query_inflight_count{host="127.0.0.1",name="postgres",port="5432",query="SELECT 1",user="postgres"} 1
			`,
			afterInflightCountExpected: `
			# HELP pg_query_inflight_count requests that have been submitted but have not been completed.
			# TYPE pg_query_inflight_count gauge
			pg_query_inflight_count{host="127.0.0.1",name="postgres",port="5432",query="SELECT 1",user="postgres"} 0
			`,
			errorCountExpected: `
			# HELP pg_query_error_count error count when querying to db
			# TYPE pg_query_error_count counter
			pg_query_error_count{host="127.0.0.1",name="postgres",port="5432",query="SELECT 1",user="postgres",} 1
			`,
		},
		"valid validation and no error": {
			err:             false,
			validValidation: true,
			query:           "SELECT 1",
			prometheusMode:  true,
			labels:          prometheus.Labels{"host": "127.0.0.1", "port": "5432", "user": "postgres", "name": "postgres"},
			beforeInflightCountExpected: `
			# HELP pg_query_inflight_count requests that have been submitted but have not been completed.
			# TYPE pg_query_inflight_count gauge
			pg_query_inflight_count{host="127.0.0.1",name="postgres",port="5432",query="SELECT 1",user="postgres"} 1
			`,
			afterInflightCountExpected: `
			# HELP pg_query_inflight_count requests that have been submitted but have not been completed.
			# TYPE pg_query_inflight_count gauge
			pg_query_inflight_count{host="127.0.0.1",name="postgres",port="5432",query="SELECT 1",user="postgres"} 0
			`,
			errorCountExpected: "",
		},
	}

	for name, s := range scenarios {
		t.Run(name, func(t *testing.T) {
			fn := func() error {
				require.Nil(t, testutil.CollectAndCompare(pgQueryInflightCount, strings.NewReader(s.beforeInflightCountExpected)))
				if s.err {
					return errors.New("")
				}

				return nil
			}

			err := DBQueryMetrics(s.labels, s.query, s.prometheusMode, fn)
			defer func() {
				pgQueryInflightCount.Reset()
				pgQueryErrorCount.Reset()
				pgQueryDurationSeconds.Reset()
			}()
			if err != nil && s.validValidation {
				require.Nil(t, testutil.CollectAndCompare(pgQueryInflightCount, strings.NewReader(s.afterInflightCountExpected)))
				require.Nil(t, testutil.CollectAndCompare(pgQueryErrorCount, strings.NewReader(s.errorCountExpected)))
				return
			} else if err != nil && !s.validValidation {
				require.Nil(t, testutil.CollectAndCompare(pgQueryInflightCount, strings.NewReader(s.afterInflightCountExpected)))
				require.Nil(t, testutil.CollectAndCompare(pgQueryErrorCount, strings.NewReader(s.errorCountExpected)))
				require.Nil(t, testutil.CollectAndCompare(pgQueryDurationSeconds, strings.NewReader(s.durationSecondsExpected)))
				return
			}

			if s.prometheusMode == false {
				require.Nil(t, testutil.CollectAndCompare(pgQueryInflightCount, strings.NewReader(s.afterInflightCountExpected)))
				require.Nil(t, testutil.CollectAndCompare(pgQueryErrorCount, strings.NewReader(s.errorCountExpected)))
				require.Nil(t, testutil.CollectAndCompare(pgQueryDurationSeconds, strings.NewReader(s.durationSecondsExpected)))
				return
			}

			require.Nil(t, testutil.CollectAndCompare(pgQueryInflightCount, strings.NewReader(s.afterInflightCountExpected)))
			require.Nil(t, testutil.CollectAndCompare(pgQueryErrorCount, strings.NewReader(s.errorCountExpected)))
		})
	}
}
