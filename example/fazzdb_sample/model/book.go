package model

import "github.com/payfazz/go-apt/pkg/fazzdb"

type Book struct {
	fazzdb.Model
	Id       string `db:"id"`
	Title    string `db:"title"`
	Year     int    `db:"year"`
	Stock    int    `db:"stock"`
	Status   string `db:"status"`
	AuthorId int    `db:"authorId"`
}

func (m *Book) Get(key string) interface{} {
	return m.Payload()[key]
}

func (m *Book) Payload() map[string]interface{} {
	return m.MapPayload(m)
}

func (m *Book) GeneratePK() {
	m.GenerateId(m)
}

func BookModel() *Book {
	return &Book{
		Model: fazzdb.UuidModel(
			"books",
			[]fazzdb.Column{
				fazzdb.Col("id"),
				fazzdb.Col("title"),
				fazzdb.Col("year"),
				fazzdb.Col("stock"),
				fazzdb.Col("status"),
				fazzdb.Col("authorId"),
			},
			"id",
			true,
			true,
		),
	}
}
