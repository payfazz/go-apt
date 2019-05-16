package qson

import "github.com/joncalhoun/qson"

// ParserInterface is a parser interface for QSON
type ParserInterface interface {
	Unmarshal(dst interface{}, query string) error
	ToJSON(query string) ([]byte, error)
}

type qs struct{}

func (q *qs) Unmarshal(dst interface{}, query string) error {
	return qson.Unmarshal(dst, query)
}

func (q *qs) ToJSON(query string) ([]byte, error) {
	return qson.ToJSON(query)
}

// NewQSONParser is a constructor that used to construct qson object
func NewQSONParser() ParserInterface {
	return &qs{}
}
