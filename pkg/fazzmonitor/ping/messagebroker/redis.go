package messagebroker

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"

	"github.com/payfazz/go-apt/pkg/fazzmonitor/ping"
)

type RedisReport struct {
	options *redis.Options
	isCore  bool
}

func (rds *RedisReport) IsCoreService() bool {
	return rds.isCore
}

func (rds *RedisReport) Check(level int64) *ping.Report {
	if !rds.isCore && level < 0 {
		return nil
	}

	report := &ping.Report{
		Service:  "redis",
		Status:   ping.NOT_AVAILABLE,
		Children: []*ping.Report{},
		IsCore:   rds.isCore,
	}

	start := time.Now()

	redisClient := redis.NewClient(rds.options)
	defer redisClient.Close()

	if err := redisClient.Ping().Err(); err != nil {
		report.Latency = ping.GetMillisecondDuration(start)
		report.Message = err.Error()

		return report
	}

	report.Latency = ping.GetMillisecondDuration(start)
	report.Status = ping.AVAILABLE

	return report
}

func NewRedisReportWithOptions(options *redis.Options, isCore bool) ping.ReportInterface {
	return &RedisReport{
		options: options,
		isCore:  isCore,
	}
}

func NewRedisReportWithAddress(hostWithPort string, password string, isCore bool) ping.ReportInterface {
	return NewRedisReportWithOptions(&redis.Options{
		Addr:     hostWithPort,
		Password: password,
		DB:       0,
	}, isCore)
}

func NewRedisReport(host string, port string, password string, isCore bool) ping.ReportInterface {
	return NewRedisReportWithAddress(fmt.Sprintf("%s:%s", host, port), password, isCore)
}
