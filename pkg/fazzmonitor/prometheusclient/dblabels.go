package prometheusclient

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

var whitelist = []string{"host", "port", "name", "user"}

func GetRequiredDBLabels() []string {
	return whitelist
}

func IsValidRequiredDBLabels(labels prometheus.Labels) bool {
	for _, name := range whitelist {
		if _, ok := labels[name]; !ok {
			log.Printf("[DBMetrics] %s is not exists & cannot continue monitor db", name)
			return false
		}
	}

	return true
}
