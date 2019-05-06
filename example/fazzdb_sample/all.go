package fazzdb_sample

import (
	"fmt"

	"github.com/payfazz/go-apt/example/fazzdb_sample/model"
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

func AllAuthors(query *fazzdb.Query) {
	conditions := []fazzdb.SliceCondition{
		{Connector: fazzdb.CO_OR, Field: "country", Operator: fazzdb.OP_EQUALS, Value: "United States"},
		{Connector: fazzdb.CO_OR, Field: "country", Operator: fazzdb.OP_EQUALS, Value: "Japan"},
		{Connector: fazzdb.CO_OR, Field: "country", Operator: fazzdb.OP_EQUALS, Value: "Singapore"},
	}

	rows, err := query.Use(model.AuthorModel()).
		//WhereIn("country", "United States", "Japan", "Singapore").
		WhereMany(conditions...).
		OrderBy("country", fazzdb.DIR_ASC).
		All()
	if nil != err {
		panic(err)
	}

	fmt.Println("--------------")
	fmt.Println("All authors from United States, Japan, Or Singapore")
	authors := rows.([]*model.Author)
	for _, author := range authors {
		fmt.Println("Name:", author.Name, ", Country:", author.Country, ", Id:", author.Id, ", CreatedAt:", author.CreatedAt)
	}
}

func AllBooks(query *fazzdb.Query) {
	rows, err := query.Use(model.BookModel()).
		WhereNil(fazzdb.UPDATED_AT).
		OrderBy("year", fazzdb.DIR_DESC).
		All()
	if nil != err {
		panic(err)
	}

	fmt.Println("--------------")
	fmt.Println("All books that is never updated")
	books := rows.([]*model.Book)
	for _, book := range books {
		fmt.Println("Title:", book.Title, ", Year:", book.Year, ", Stock:", book.Stock, ", CreatedAt:", book.CreatedAt)
	}
}
