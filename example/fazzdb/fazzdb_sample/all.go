package fazzdb_sample

import (
	"fmt"
	model2 "github.com/payfazz/go-apt/example/fazzdb/fazzdb_sample/model"

	"github.com/payfazz/go-apt/pkg/fazzdb"
)

func AllAuthors(query *fazzdb.Query) {
	conditions := []fazzdb.SliceCondition{
		{Connector: fazzdb.CO_OR, Field: "country", Operator: fazzdb.OP_EQUALS, Value: "United States"},
		{Connector: fazzdb.CO_OR, Field: "country", Operator: fazzdb.OP_EQUALS, Value: "Japan"},
		{Connector: fazzdb.CO_OR, Field: "country", Operator: fazzdb.OP_EQUALS, Value: "Singapore"},
	}

	rows, err := query.Use(model2.AuthorModel()).
		//WhereIn("country", "United States", "Japan", "Singapore").
		WhereMany(conditions...).
		OrderBy("country", fazzdb.DIR_ASC).
		All()
	if nil != err {
		panic(err)
	}

	fmt.Println("--------------")
	fmt.Println("All authors from United States, Japan, Or Singapore")
	authors := rows.([]*model2.Author)
	for _, author := range authors {
		fmt.Println("Name:", author.Name, ", Country:", author.Country, ", Id:", author.Id, ", CreatedAt:", author.CreatedAt)
	}
}

func AllBooks(query *fazzdb.Query) {
	rows, err := query.Use(model2.BookModel()).
		WhereNil(fazzdb.UPDATED_AT).
		OrderBy("year", fazzdb.DIR_DESC).
		All()
	if nil != err {
		panic(err)
	}

	fmt.Println("--------------")
	fmt.Println("All books that is never updated")
	books := rows.([]*model2.Book)
	for _, book := range books {
		fmt.Println("Title:", book.Title, ", Year:", book.Year, ", Stock:", book.Stock, ", CreatedAt:", book.CreatedAt)
	}
}

func AllBooksSliceConditions(query *fazzdb.Query) {
	conditions := []fazzdb.SliceCondition{
		{
			Connector: fazzdb.CO_OR,
			Conditions: []fazzdb.SliceCondition{
				{Connector: fazzdb.CO_OR, Field: "year", Operator: fazzdb.OP_EQUALS, Value: 1967},
				{Connector: fazzdb.CO_OR, Field: "year", Operator: fazzdb.OP_EQUALS, Value: 1972},
			},
		},
		{Connector: fazzdb.CO_OR, Field: "year", Operator: fazzdb.OP_EQUALS, Value: 1995},
		{
			Connector: fazzdb.CO_OR,
			Conditions: []fazzdb.SliceCondition{
				{Connector: fazzdb.CO_OR, Field: "year", Operator: fazzdb.OP_EQUALS, Value: 1978},
				{Connector: fazzdb.CO_OR, Field: "year", Operator: fazzdb.OP_EQUALS, Value: 1982},
				{Connector: fazzdb.CO_OR, Field: "year", Operator: fazzdb.OP_EQUALS, Value: 1983},
				{Connector: fazzdb.CO_OR, Field: "year", Operator: fazzdb.OP_EQUALS, Value: 1984},
			},
		},
		{Connector: fazzdb.CO_OR, Field: "year", Operator: fazzdb.OP_EQUALS, Value: 1994},
		{
			Connector: fazzdb.CO_OR,
			Conditions: []fazzdb.SliceCondition{
				{Connector: fazzdb.CO_OR, Field: "year", Operator: fazzdb.OP_EQUALS, Value: 1989},
				{Connector: fazzdb.CO_OR, Field: "year", Operator: fazzdb.OP_EQUALS, Value: 1990},
				{Connector: fazzdb.CO_OR, Field: "year", Operator: fazzdb.OP_EQUALS, Value: 1988},
				{Connector: fazzdb.CO_OR, Field: "year", Operator: fazzdb.OP_EQUALS, Value: 1987},
				{Connector: fazzdb.CO_OR, Field: "year", Operator: fazzdb.OP_EQUALS, Value: 1986},
				{Connector: fazzdb.CO_OR, Field: "year", Operator: fazzdb.OP_EQUALS, Value: 1985},
			},
		},
		{Connector: fazzdb.CO_OR, Field: "year", Operator: fazzdb.OP_EQUALS, Value: 1993},
		{Connector: fazzdb.CO_OR, Field: "year", Operator: fazzdb.OP_EQUALS, Value: 1992},
		{Connector: fazzdb.CO_OR, Field: "year", Operator: fazzdb.OP_EQUALS, Value: 1991},
		{Connector: fazzdb.CO_OR, Field: "year", Operator: fazzdb.OP_EQUALS, Value: 1990},
	}

	rows, err := query.Use(model2.BookModel()).
		WhereMany(conditions...).
		All()
	if nil != err {
		panic(err)
	}

	fmt.Println("--------------")
	fmt.Println("Test prefix using group slice conditions")
	books := rows.([]*model2.Book)
	for _, book := range books {
		fmt.Println("Title:", book.Title, ", Year:", book.Year, ", Stock:", book.Stock, ", CreatedAt:", book.CreatedAt)
	}
}
