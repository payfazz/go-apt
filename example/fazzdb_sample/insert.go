package fazzdb_sample

import (
	"github.com/payfazz/go-apt/example/fazzdb_sample/model"
	"github.com/payfazz/go-apt/pkg/fazzdb"
	"log"
)

func InsertBook(query *fazzdb.Query) {
	row, err := query.Use(model.AuthorModel()).First()
	if nil != err {
		log.Fatalln(err)
	}
	author := row.(*model.Author)

	book := model.BookModel()
	book.Title = "Artemis Fowl"
	book.Year = 2010
	book.Stock = 20
	book.Status = model.BOOK_AVAILABLE
	book.AuthorId = author.Id

	_, err = query.Use(book).Insert()
	if nil != err {
		panic(err)
	}
}

func InsertAuthor(query *fazzdb.Query) {
	author := model.AuthorModel()
	author.Name = "Eoin Colfer"
	author.Country = "United Kingdom"

	_, err := query.Use(author).Insert()
	if nil != err {
		panic(err)
	}
}
