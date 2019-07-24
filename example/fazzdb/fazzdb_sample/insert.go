package fazzdb_sample

import (
	model2 "github.com/payfazz/go-apt/example/fazzdb/fazzdb_sample/model"
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

func InsertBook(query *fazzdb.Query) {
	row, err := query.Use(model2.AuthorModel()).First()
	if nil != err {
		panic(err)
	}
	author := row.(*model2.Author)

	book := model2.BookModel()
	book.Title = "Artemis Fowl"
	book.Year = 2010
	book.Stock = 20
	book.Status = model2.BOOK_AVAILABLE
	book.AuthorId = author.Id

	_, err = query.Use(book).Insert()
	if nil != err {
		panic(err)
	}
}

func InsertAuthor(query *fazzdb.Query) {
	author := model2.AuthorModel()
	author.Name = "Eoin Colfer"
	author.Country = "United Kingdom"

	_, err := query.Use(author).Insert()
	if nil != err {
		panic(err)
	}
}
