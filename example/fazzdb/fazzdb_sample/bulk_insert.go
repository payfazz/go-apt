package fazzdb_sample

import (
	"fmt"
	model2 "github.com/payfazz/go-apt/example/fazzdb/fazzdb_sample/model"
	"github.com/payfazz/go-apt/pkg/fazzdb"
	"math/rand"
)

var authorNames = []string{
	"Jane Austen",
	"Charles Dickens",
	"Ernest Hemingway",
	"Fitzgerald",
	"Dan Brown",
	"Harper Lee",
	"Joseph Conrad",
	"William Faulkner",
	"Arthur Conan Doyle",
	"Agatha Christie",
}
var countryNames = []string{
	"United States",
	"Japan",
	"Singapore",
	"Hong Kong",
	"German",
}

func BulkInsertAuthors(query *fazzdb.Query) {
	var authors []*model2.Author
	for i := 0; i < 10; i++ {
		newAuthor := model2.AuthorModel()
		newAuthor.Name = authorNames[rand.Intn(len(authorNames))]
		newAuthor.Country = countryNames[rand.Intn(len(countryNames))]

		authors = append(authors, newAuthor)
	}

	_, err := query.Use(model2.AuthorModel()).BulkInsert(authors)
	if nil != err {
		panic(err)
	}
}

func BulkInsertBooks(query *fazzdb.Query) {
	var books []*model2.Book

	count, err := query.Use(model2.AuthorModel()).Count()
	if nil != err {
		panic(err)
	}

	authorCount := int(*count) - 2

	for i := 0; i < 10; i++ {
		newBook := model2.BookModel()
		newBook.Title = fmt.Sprintf("%s Books", authorNames[rand.Intn(len(authorNames))])
		newBook.Year = rand.Intn(100) + 1920
		newBook.Stock = rand.Intn(30) + 10
		newBook.Status = model2.BOOK_BORROWED
		newBook.AuthorId = rand.Intn(authorCount) + 1

		books = append(books, newBook)
	}

	_, err = query.Use(model2.BookModel()).BulkInsert(books)
	if nil != err {
		panic(err)
	}
}
