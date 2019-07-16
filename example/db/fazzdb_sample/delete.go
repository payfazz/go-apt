package fazzdb_sample

import (
	"fmt"
	model2 "github.com/payfazz/go-apt/example/db/fazzdb_sample/model"
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

func DeleteAuthor(query *fazzdb.Query) {
	row, err := query.Use(model2.AuthorModel()).
		Where("country", "United States").
		First()
	if nil != err {
		panic(err)
	}

	author := row.(*model2.Author)
	fmt.Println("--------------")
	fmt.Println("Removing first Author from United States:", author.Name, ", Id:", author.Id)
	_, err = query.Use(author).Delete()
	if nil != err {
		panic(err)
	}
}

func DeleteBook(query *fazzdb.Query) {
	row, err := query.Use(model2.BookModel()).
		WhereOp("year", fazzdb.OP_LESS_THAN, 1980).
		First()
	if nil != err {
		panic(err)
	}

	book := row.(*model2.Book)
	fmt.Println("--------------")
	fmt.Println("Removing first Book before 1980:", book.Title, ", Stock:", book.Stock, ", Id:", book.Id)
	_, err = query.Use(book).Delete()
	if nil != err {
		panic(err)
	}
}
