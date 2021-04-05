package filter

import (
	"math"
	"net/url"
	"reflect"
	"time"

	"github.com/jinzhu/now"

	"github.com/payfazz/go-apt/pkg/fazzcommon/formatter"
	"github.com/payfazz/go-apt/pkg/fazzcommon/value/pagination"
	"github.com/payfazz/go-apt/pkg/fazzcommon/value/timestamp"
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
	limit := defaultLimit
	page := 1

	if limitQuery := queryParams.Get(pagination.LIMIT_PARAM); "" != limitQuery {
		limit = formatter.StringToInteger(limitQuery)
	}

	if pageQuery := queryParams.Get(pagination.PAGE_PARAM); "" != pageQuery {
		page = int(math.Max(float64(1), formatter.StringToFloat(pageQuery)))
	}

	return BuildPage(limit, page)
}

// BuildPage is a function to generate Page based on given limit and page
func BuildPage(limit int, page int) *Page {
	finalLimit := limit + 1
	if limit == -1 {
		finalLimit = limit
	}

	return &Page{
		BaseLimit: limit,
		Limit:     finalLimit,
		Page:      page,
		Offset:    (page - 1) * limit,
	}
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

	if page.BaseLimit == -1 {
		hasNext = false
	}

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
	return ParseTimestampRangeInLocation(queryParams, defaultStart, time.UTC)
}

// ParseTimestampRangeInLocation is a function to handle start and end date payload in location
func ParseTimestampRangeInLocation(queryParams url.Values, defaultStart *time.Time, loc *time.Location) (*TimestampRange, error) {
	currentTime := time.Now()
	timestampRange := &TimestampRange{
		Start: defaultStart,
		End:   &currentTime,
	}

	if start := queryParams.Get(timestamp.START_PARAM); "" != start {
		startTime, err := ParseTimestampInLocation(start, loc)
		if nil != err {
			return nil, err
		}

		timestampRange.Start = startTime
	}

	if end := queryParams.Get(timestamp.END_PARAM); "" != end {
		endTime, err := ParseTimestampInLocation(end, loc)
		if nil != err {
			return nil, err
		}

		// set endTime as end of day if time set as 00:00:00
		endTimeEndOfDay := *endTime
		if endTime.Hour() != 0 || endTime.Minute() == 0 && endTime.Second() == 0 {
			endTimeEndOfDay = now.With(endTimeEndOfDay).EndOfDay()
		}
		timestampRange.End = &endTimeEndOfDay
	}

	return timestampRange, nil
}

// ParseTimestamp is a function to handle converting multiple Date / Datetime format to *time.Time
func ParseTimestamp(arg string) (*time.Time, error) {
	return ParseTimestampInLocation(arg, time.UTC)
}

// ParseTimestampInLocation is a function to handle converting multiple Date / Datetime format to *time.Time in location
func ParseTimestampInLocation(arg string, loc *time.Location) (*time.Time, error) {
	formats := []string{
		timestamp.TS_RFC3339,
		timestamp.TS_DATE,
		timestamp.TS_DATETIME,
		timestamp.TS_DATETIME_WITH_MILISECONDS,
	}

	for _, f := range formats {
		t, err := time.ParseInLocation(f, arg, loc)
		if nil == err {
			return &t, nil
		}
	}

	return nil, errors.New("Unsupported timestamp format")
}
