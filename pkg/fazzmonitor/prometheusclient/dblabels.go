package prometheusclient

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

var labels = []string{"host", "port", "name", "user"}

// GetRequiredDBLabels return the required labels for infra
func GetRequiredDBLabels() []string {
	return labels
}

// IsValidRequiredDBLabels check if the argument labels is satify infra or not
func IsValidRequiredDBLabels(l prometheus.Labels) bool {
	valid := true
	for _, name := range labels {
		if _, ok := l[name]; !ok {
			log.Printf("[DBMetrics] %s is not exists & cannot continue monitor db", name)
			valid = false
		}
	}

	return valid
}
