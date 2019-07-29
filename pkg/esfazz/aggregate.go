package esfazz

import (
	"encoding/json"
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

// AggregateRow is a model for aggregate snapshot in database
type AggregateRow struct {
	fazzdb.Model
	Id      string          `json:"id" db:"id"`
	Version int             `json:"version" db:"version"`
	Data    json.RawMessage `json:"data" db:"data"`
}

// GeneratePK is a function used by fazzdb to generate PK
func (m *AggregateRow) GeneratePK() {
	// PK manually inserted
}

// Get is a function that used to get the data from table
func (m *AggregateRow) Get(key string) interface{} {
	return m.Payload()[key]
}

// Payload is a function that used to get the payload data
func (m *AggregateRow) Payload() map[string]interface{} {
	return m.MapPayload(m)
}

func AggregateRowModel(table string) *AggregateRow {
	return &AggregateRow{
		Model: fazzdb.UuidModel(table,
			[]fazzdb.Column{
				fazzdb.Col("id"),
				fazzdb.Col("version"),
				fazzdb.Col("data"),
			},
			"id",
			false,
			false,
		),
	}
}

// Aggregate is interface for aggregate object
type Aggregate interface {
	GetId() string
	GetVersion() int
}

// BaseAggregate is a struct to be used composed with aggregate object
type BaseAggregate struct {
	Id      string `json:"id"`
	Version int    `json:"version"`
}

// GetId return Id of aggregate object
func (a *BaseAggregate) GetId() string {
	return a.Id
}

// GetVersion return aggregate version of aggregate object
func (a *BaseAggregate) GetVersion() int {
	return a.Version
}
