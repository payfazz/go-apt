package database

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/payfazz/go-apt/pkg/fazzmonitor/ping"
)

type MySQLReport struct {
	connectionString string
}

func (mysql *MySQLReport) Check(level int64) ping.Report {
	report := ping.Report{
		Service:  "mysql",
		Status:   ping.NOT_AVAILABLE,
		Children: []ping.Report{},
	}

	start := time.Now()

	db, err := sqlx.Connect("mysql", mysql.connectionString)
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

func NewMySQLReportWithConnectionString(connectionString string) ping.ReportInterface {
	return &MySQLReport{
		connectionString: connectionString,
	}
}

func NewMySQLReport(host string, port string, user string, password string, dbName string) ping.ReportInterface {
	return NewMySQLReportWithConnectionString(fmt.Sprintf(
		"%s:%s@(%s:%s)/%s",
		user, password, host, port, dbName,
	))
}
