package filter

import (
	"math"
	"net/url"
	"reflect"
	"time"

	"github.com/payfazz/fazzcard-backend/lib/fazzcommon/value/pagination"
	"github.com/payfazz/fazzcard-backend/lib/fazzcommon/value/timestamp"
	"github.com/payfazz/go-apt/pkg/fazzcommon/formatter"
	"github.com/pkg/errors"
)

// Page is a struct to handle paging attributes
type Page struct {
	Limit     int
	BaseLimit int `json:"limit"`
	Page      int `json:"page"`
	Offset    int `json:"offset"`
}

// ParsePage is a function to parse paging attribute from query params to Page struct
func ParsePage(queryParams url.Values, defaultLimit int) *Page {
	p := &Page{
		BaseLimit: defaultLimit,
		Limit:     defaultLimit + 1,
		Page:      1,
	}

	if limit := queryParams.Get(pagination.LIMIT_PARAM); "" != limit {
		lm := formatter.StringToInteger(limit)
		p.BaseLimit = lm
		p.Limit = lm + 1
	}

	if page := queryParams.Get(pagination.PAGE_PARAM); "" != page {
		p.Page = int(math.Max(float64(1), formatter.StringToFloat(page)))
	}

	p.Offset = (p.Page - 1) * p.BaseLimit

	return p
}

// PageResponse is a struct to handle data with paging details
type PageResponse struct {
	Data    interface{} `json:"data"`
	Count   int         `json:"count"`
	HasNext *bool       `json:"hasNext"`
}

// BuildPageResponse is a function to build response with additional page data
func BuildPageResponse(page *Page, data interface{}) (*PageResponse, error) {
	if reflect.Slice != reflect.TypeOf(data).Kind() {
		return nil, errors.New("paginated data must be a slice")
	}

	s := reflect.ValueOf(data)
	hasNext := s.Len() > page.BaseLimit
	length := s.Len()

	result := s.Interface()
	if hasNext {
		result = s.Slice(0, s.Len()-1).Interface()
		length = s.Len() - 1
	}

	if s.Len() == 0 {
		result = []struct{}{}
	}

	return &PageResponse{
		Data:    result,
		Count:   length,
		HasNext: &hasNext,
	}, nil
}

// TimestampRange is a struct to handle date range attributes
type TimestampRange struct {
	Start *time.Time `json:"start"`
	End   *time.Time `json:"end"`
}

// ParseTimestampRange is a function to handle start and end date payload
func ParseTimestampRange(queryParams url.Values, defaultStart *time.Time) (*TimestampRange, error) {
	now := time.Now()
	timestampRange := &TimestampRange{
		Start: defaultStart,
		End:   &now,
	}

	if start := queryParams.Get(timestamp.START_PARAM); "" != start {
		startTime, err := ParseTimestamp(start)
		if nil != err {
			return nil, err
		}

		timestampRange.Start = startTime
	}

	if end := queryParams.Get(timestamp.END_PARAM); "" != end {
		endTime, err := ParseTimestamp(end)
		if nil != err {
			return nil, err
		}

		timestampRange.End = endTime
	}

	return timestampRange, nil
}

// ParseTimestamp is a function to handle converting multiple Date / Datetime format to *time.Time
func ParseTimestamp(arg string) (*time.Time, error) {
	formats := []string{
		timestamp.TS_RFC3339,
		timestamp.TS_DATE,
		timestamp.TS_DATETIME,
		timestamp.TS_DATETIME_WITH_MILISECONDS,
	}

	for _, f := range formats {
		t, err := time.Parse(f, arg)
		if nil == err {
			return &t, nil
		}
	}

	return nil, errors.New("Unsupported timestamp format")
}
