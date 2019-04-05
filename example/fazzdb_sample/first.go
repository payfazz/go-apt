package fazzdb_sample

import (
	"fmt"
	"github.com/payfazz/go-apt/example/fazzdb_sample/model"
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

func FirstAuthor(query *fazzdb.Query) {
	row, err := query.Use(model.AuthorModel()).
		Where("country", "United States").
		First()
	if nil != err {
		panic(err)
	}

	author := row.(*model.Author)
	fmt.Println("--------------")
	fmt.Println("First Author from United States:", author.Name, ", Id:", author.Id)
}

func FirstBook(query *fazzdb.Query) {
	row, err := query.Use(model.BookModel()).
		WhereOp("year", fazzdb.OP_LESS_THAN, 1960).
		OrWhereOp("year", fazzdb.OP_MORE_THAN, 1990).
		First()
	if nil != err {
		panic(err)
	}

	book := row.(*model.Book)
	fmt.Println("--------------")
	fmt.Println("First Book between 1960 and 1990:", book.Title, ", Stock:", book.Stock, ", Id:", book.Id)
}
