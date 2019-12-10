package prometheusclient

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var dateMinuteFormat = "2006-01-02 15:04"

func registerOnce(collector prometheus.Collector) prometheus.Collector {
	if err := prometheus.DefaultRegisterer.Register(collector); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			return are.ExistingCollector
		} else {
			panic(err)
		}
	}

	return collector
}

func dateMinute() string {
	return time.Now().Format(dateMinuteFormat)
}

func path(r *http.Request) string {
	return fmt.Sprintf("%s%s", r.Host, r.URL.Path)
}
