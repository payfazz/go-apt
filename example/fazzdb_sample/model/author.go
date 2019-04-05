package model

import "github.com/payfazz/go-apt/pkg/fazzdb"

type Author struct {
	fazzdb.Model
	Id      int    `db:"id"`
	Name    string `db:"name"`
	Country string `db:"country"`
}

func (m *Author) Get(key string) interface{} {
	return m.Payload()[key]
}

func (m *Author) Payload() map[string]interface{} {
	return m.MapPayload(m)
}

func AuthorModel() *Author {
	return &Author{
		Model: fazzdb.AutoIncrementModel(
			"authors",
			[]fazzdb.Column{
				fazzdb.Col("id"),
				fazzdb.Col("name"),
				fazzdb.Col("country"),
			},
			"id",
			true,
			false,
		),
	}
}