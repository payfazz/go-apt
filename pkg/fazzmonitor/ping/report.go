package ping

import (
	"encoding/json"
	"net/http"
	"time"
)

type ReportInterface interface {
	Check() Report
}

type Report struct {
	Latency     int64  `json:"latency"`
	IsAvailable bool   `json:"isAvailable"`
	Message     string `json:"message"`
}

func GetMillisecondDuration(startRequestAt time.Time) int64 {
	return time.Since(startRequestAt).Milliseconds()
}

func Ping(reportChecks map[string]ReportInterface) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(writer http.ResponseWriter, request *http.Request) {
			result := make(map[string]Report, 0)

			for component, reportCheck := range reportChecks {
				result[component] = reportCheck.Check()
			}

			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(writer).Encode(result)
			return
		}
	}
}

func PingHandler(reportChecks map[string]ReportInterface) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return Ping(reportChecks)(next.ServeHTTP)
	}
}
