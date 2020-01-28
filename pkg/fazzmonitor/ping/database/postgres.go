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
}

func (pg *PgSQLReport) Check() ping.Report {
	report := ping.Report{
		IsAvailable: false,
	}

	start := time.Now()

	db, err := sqlx.Connect("postgres", pg.connectionString)
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
	report.IsAvailable = true

	return report
}

func NewPgSQLReportWithConnectionString(connectionString string) ping.ReportInterface {
	return &PgSQLReport{
		connectionString: connectionString,
	}
}

func NewPgSQLReport(host string, port string, user string, password string, dbName string) ping.ReportInterface {
	return NewPgSQLReportWithConnectionString(fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName,
	))
}
