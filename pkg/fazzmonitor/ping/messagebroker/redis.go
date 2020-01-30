package messagebroker

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"

	"github.com/payfazz/go-apt/pkg/fazzmonitor/ping"
)

type RedisReport struct {
	options *redis.Options
}

func (rds *RedisReport) Check(level int64) ping.Report {
	report := ping.Report{
		Service: "redis",
		Status:  ping.NOT_AVAILABLE,
	}

	start := time.Now()

	redisClient := redis.NewClient(rds.options)

	if err := redisClient.Ping().Err(); err != nil {
		report.Latency = ping.GetMillisecondDuration(start)
		report.Message = err.Error()

		return report
	}

	report.Latency = ping.GetMillisecondDuration(start)
	report.Status = ping.AVAILABLE

	return report
}

func NewRedisReportWithOptions(options *redis.Options) ping.ReportInterface {
	return &RedisReport{
		options: options,
	}
}

func NewRedisReportWithAddress(hostWithPort string, password string) ping.ReportInterface {
	return NewRedisReportWithOptions(&redis.Options{
		Addr:     hostWithPort,
		Password: password,
		DB:       0,
	})
}

func NewRedisReport(host string, port string, password string) ping.ReportInterface {
	return NewRedisReportWithAddress(fmt.Sprintf("%s:%s", host, port), password)
}
