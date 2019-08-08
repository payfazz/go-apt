package fazzdb_sample

import (
	"fmt"
	"github.com/payfazz/go-apt/example/fazzdb/fazzdb_sample/model"
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

func DeleteAuthor(query *fazzdb.Query) {
	row, err := query.Use(model.AuthorModel()).
		Where("country", "United States").
		First()
	if nil != err {
		panic(err)
	}

	author := row.(*model.Author)
	fmt.Println("--------------")
	fmt.Println("Removing first Author from United States:", author.Name, ", Id:", author.Id)
	_, err = query.Use(author).Delete()
	if nil != err {
		panic(err)
	}
}

func DeleteBook(query *fazzdb.Query) {
	row, err := query.Use(model.BookModel()).
		WhereOp("year", fazzdb.OP_LESS_THAN, 1980).
		First()
	if nil != err {
		panic(err)
	}

	book := row.(*model.Book)
	fmt.Println("--------------")
	fmt.Println("Removing first Book before 1980:", book.Title, ", Stock:", book.Stock, ", Id:", book.Id)
	_, err = query.Use(book).Delete()
	if nil != err {
		panic(err)
	}
}
