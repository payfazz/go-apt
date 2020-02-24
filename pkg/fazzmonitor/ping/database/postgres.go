package database

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/payfazz/go-apt/pkg/fazzmonitor/ping"
)

type PgSQLReport struct {
	connectionString string
	isCore           bool
}

func (pg *PgSQLReport) IsCoreService() bool {
	return pg.isCore
}

func (pg *PgSQLReport) Check(level int64) *ping.Report {
	if !pg.isCore && level < 0 {
		return nil
	}

	report := &ping.Report{
		Service:  "postgres",
		Status:   ping.NOT_AVAILABLE,
		Children: []*ping.Report{},
		IsCore:   pg.isCore,
	}

	start := time.Now()

	db, err := sqlx.Connect("postgres", pg.connectionString)
	defer db.Close()
	if nil != err {
		report.Latency = ping.GetMillisecondDuration(start)
		report.Message = err.Error()

		return report
	}

	_, err = db.Exec("SELECT 1;")
	if nil != err {
		report.Latency = ping.GetMillisecondDuration(start)
		report.Message = err.Error()

		return report
	}

	report.Latency = ping.GetMillisecondDuration(start)
	report.Status = ping.AVAILABLE

	return report
}

func NewPgSQLReportWithConnectionString(connectionString string, isCore bool) ping.ReportInterface {
	return &PgSQLReport{
		connectionString: connectionString,
		isCore:           isCore,
	}
}

func NewPgSQLReport(host string, port string, user string, password string, dbName string, isCore bool) ping.ReportInterface {
	return NewPgSQLReportWithConnectionString(fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName,
	), isCore)
}
