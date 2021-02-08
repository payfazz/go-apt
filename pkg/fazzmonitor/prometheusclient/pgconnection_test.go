package prometheusclient

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/require"
)

func TestPGConnection(t *testing.T) {
	host := os.Getenv("TEST_DB_HOST")
	port := os.Getenv("TEST_DB_PORT")
	user := os.Getenv("TEST_DB_USER")
	pass := os.Getenv("TEST_DB_PASS")
	name := os.Getenv("TEST_DB_NAME")
	scenarios := map[string]struct {
		validValidation        bool
		labels                 prometheus.Labels
		idleConnectionExpected string
		useConnectionExpected  string
		waitConnectionExpected string
	}{
		"invalid validation": {
			validValidation: false,
		},
		"valid validation": {
			validValidation: true,
			labels: prometheus.Labels{
				"host": host,
				"port": port,
				"user": user,
				"name": name,
			},
			idleConnectionExpected: fmt.Sprintf(`
				# HELP pg_connection_idle_count show the database connection idle
				# TYPE pg_connection_idle_count gauge
				pg_connection_idle_count{host="%s",name="%s",port="%s",user="%s"} 0
			`, host, name, port, user),
			useConnectionExpected: fmt.Sprintf(`
				# HELP pg_connection_use_count show the database connection use count
				# TYPE pg_connection_use_count gauge
				pg_connection_use_count{host="%s",name="%s",port="%s",user="%s"} 0
			`, host, name, port, user),
			waitConnectionExpected: fmt.Sprintf(`
				# HELP pg_connection_wait_count show the database connection wait count
				# TYPE pg_connection_wait_count gauge
				pg_connection_wait_count{host="%s",name="%s",port="%s",user="%s"} 0
			`, host, name, port, user),
		},
	}

	for name, s := range scenarios {
		t.Run(name, func(t *testing.T) {
			conn := fmt.Sprintf(
				"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
				host, port, user, pass, name,
			)
			db, err := sqlx.Open("postgres", conn)
			if err != nil {
				t.Log("Skip doing test cause the env is not set")
				return
			}
			defer db.Close()

			PGConnectionGauge(s.labels, db)
			defer func() {
				pgConnectionIdleCount.Reset()
				pgConnectionUseCount.Reset()
				pgConnectionWaitCount.Reset()
			}()
			if !s.validValidation {
				require.Nil(t, testutil.CollectAndCompare(pgConnectionIdleCount, strings.NewReader(s.idleConnectionExpected)))
				require.Nil(t, testutil.CollectAndCompare(pgConnectionUseCount, strings.NewReader(s.useConnectionExpected)))
				require.Nil(t, testutil.CollectAndCompare(pgConnectionWaitCount, strings.NewReader(s.waitConnectionExpected)))
				return
			}

			require.Nil(t, testutil.CollectAndCompare(pgConnectionIdleCount, strings.NewReader(s.idleConnectionExpected)))
			require.Nil(t, testutil.CollectAndCompare(pgConnectionUseCount, strings.NewReader(s.useConnectionExpected)))
			require.Nil(t, testutil.CollectAndCompare(pgConnectionWaitCount, strings.NewReader(s.waitConnectionExpected)))
		})
	}
}
