package query

import "github.com/payfazz/go-apt/pkg/fazzdb"

// Todo is a read model that is returned for todo service
type Todo struct {
	fazzdb.Model `json:"-"`
	Id           string `json:"id" db:"id"`
	Text         string `json:"text" db:"text"`
	Completed    bool   `json:"completed" db:"completed"`
}

// GeneratePK used to generate PK for fazzcard
func (m *Todo) GeneratePK() {
	// Don't generate primary key, use key from event
}

// Get used to get the data
func (m *Todo) Get(key string) interface{} {
	return m.Payload()[key]
}

// Payload used to get the payload
func (m *Todo) Payload() map[string]interface{} {
	return m.MapPayload(m)
}

// TodoReadModel is a constructor for user model
func TodoReadModel() *Todo {
	return &Todo{
		Model: fazzdb.UuidModel(
			"todos_read",
			[]fazzdb.Column{
				fazzdb.Col("id"),
				fazzdb.Col("text"),
				fazzdb.Col("completed"),
			},
			"id",
			false,
			false,
		),
	}
}
