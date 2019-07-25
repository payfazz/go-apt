package timestamp

import "time"

const (
	// START_PARAM start
	START_PARAM = "start"
	// END_PARAM end
	END_PARAM = "end"
)

const (
	// TS_RFC3339 time.RFC3339
	TS_RFC3339 = time.RFC3339
	// TS_DATE yyyy-mm-dd
	TS_DATE = "2006-01-02"
	// TS_DATETIME yyyy-mm-dd hh:mm:ss
	TS_DATETIME = "2006-01-02 15:04:05"
	// TS_DATETIME_WITH_MILISECONDS yyyy-mm-dd hh:mm:ss.SSS
	TS_DATETIME_WITH_MILISECONDS = "2006-01-02 15:04:05.000"
	// TS_TIME hh:mm:ss
	TS_TIME = "15:04:05"
	// TS_TIME_WITH_MILISECONDS hh:mm:ss.SSS
	TS_TIME_WITH_MILISECONDS = "15:04:05.000"
)
