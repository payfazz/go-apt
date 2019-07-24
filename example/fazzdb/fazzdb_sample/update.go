package fazzdb_sample

import (
	model2 "github.com/payfazz/go-apt/example/fazzdb/fazzdb_sample/model"
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

func UpdateAuthor(query *fazzdb.Query) {
	row, err := query.Use(model2.AuthorModel()).
		Where("id", 1).
		First()
	if nil != err {
		panic(err)
	}

	author := row.(*model2.Author)
	author.Name = "J.K Rowling"
	_, err = query.Use(author).Update()
	if nil != err {
		panic(err)
	}
}

func UpdateBook(query *fazzdb.Query) {
	row, err := query.Use(model2.BookModel()).
		First()
	if nil != err {
		panic(err)
	}

	book := row.(*model2.Book)
	book.Title = "Harry Potter"
	_, err = query.Use(book).Update()
	if nil != err {
		panic(err)
	}
}
