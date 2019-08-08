package snappostgres

import (
	"github.com/jmoiron/sqlx/types"
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

// AggregateRow is a model for aggregate snapshot in database
type AggregateRow struct {
	fazzdb.Model
	Id   string         `json:"id" db:"id"`
	Data types.JSONText `json:"data" db:"data"`
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

// AggregateRowModel is the constructor for aggregate row model
func AggregateRowModel(table string) *AggregateRow {
	return &AggregateRow{
		Model: fazzdb.UuidModel(table,
			[]fazzdb.Column{
				fazzdb.Col("id"),
				fazzdb.Col("data"),
			},
			"id",
			false,
			false,
		),
	}
}
