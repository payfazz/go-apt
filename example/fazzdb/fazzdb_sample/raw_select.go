package fazzdb_sample

import (
	"fmt"
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

type BookCompact struct {
	Id     string `db:"id"`
	Title  string `db:"title"`
	Year   int    `db:"year"`
	Status string `db:"status"`
}

func RawFirst(query *fazzdb.Query) {
	row, err := query.RawFirst(&BookCompact{}, "SELECT id, title, year, status FROM books WHERE year < $1;", 1955)
	if nil != err {
		panic(err)
	}

	book := row.(*BookCompact)
	fmt.Println("--------------")
	fmt.Println("Raw First book without author before 1955, Title:", book.Title, ",Year:", book.Year, ",Status:", book.Status)
}

func RawAll(query *fazzdb.Query) {
	rows, err := query.RawAll(&BookCompact{}, "SELECT id, title, year, status FROM books WHERE year < $1;", 1955)
	if nil != err {
		panic(err)
	}

	books := rows.([]*BookCompact)
	fmt.Println("--------------")
	fmt.Println("Raw All book without author")
	for _, v := range books {
		fmt.Println("Title:", v.Title, ", Year:", v.Year, ", Status:", v.Status)
	}
}

func RawNamedFirst(query *fazzdb.Query) {
	payload := map[string]interface{}{
		"year": 1955,
	}
	row, err := query.RawNamedFirst(&BookCompact{}, "SELECT id, title, year, status FROM books WHERE year < :year;", payload)
	if nil != err {
		panic(err)
	}

	book := row.(*BookCompact)
	fmt.Println("--------------")
	fmt.Println("Raw First book without author before 1955, Title:", book.Title, ", Year:", book.Year, ", Status:", book.Status)
}

func RawNamedAll(query *fazzdb.Query) {
	payload := map[string]interface{}{
		"year": 1955,
	}
	rows, err := query.RawNamedAll(&BookCompact{}, "SELECT id, title, year, status FROM books WHERE year < :year;", payload)
	if nil != err {
		panic(err)
	}

	books := rows.([]*BookCompact)
	fmt.Println("--------------")
	fmt.Println("Raw All book without author")
	for _, v := range books {
		fmt.Println("Title:", v.Title, ", Year:", v.Year, ", Status:", v.Status)
	}
}
