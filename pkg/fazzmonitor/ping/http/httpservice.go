package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/payfazz/go-apt/pkg/fazzmonitor/ping"
)

type HttpServiceReport struct {
	url    string
	isCore bool
}

func (s *HttpServiceReport) IsCoreService() bool {
	return s.isCore
}

func (s *HttpServiceReport) Check(level int64) *ping.Report {
	if !s.isCore && level < 0 {
		return nil
	}

	urlWithLevel := fmt.Sprintf("%s?%s=%d", s.url, ping.LEVEL_KEY, level)

	report := &ping.Report{
		Service:  s.url,
		Status:   ping.NOT_AVAILABLE,
		Children: []*ping.Report{},
		IsCore:   s.isCore,
	}

	start := time.Now()

	resp, err := http.Get(urlWithLevel)
	if nil != err {
		report.Latency = ping.GetMillisecondDuration(start)
		report.Message = err.Error()

		return report
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		report.Latency = ping.GetMillisecondDuration(start)
		report.Message = err.Error()

		return report
	}

	var serviceReport ping.Report
	err = json.Unmarshal(body, &serviceReport)
	if nil != err {
		report.Latency = ping.GetMillisecondDuration(start)
		report.Status = ping.AVAILABLE
		report.Message = fmt.Sprintf("ping report not implemented; body: %s", string(body))

		return report
	}

	return &serviceReport
}

func NewHttpServiceReport(url string, isCore bool) ping.ReportInterface {
	return &HttpServiceReport{
		url:    url,
		isCore: isCore,
	}
}
