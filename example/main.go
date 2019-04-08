package main

import (
	_ "github.com/lib/pq"
	"github.com/payfazz/go-apt/config"
	"github.com/payfazz/go-apt/example/fazzdb_sample"
	"github.com/payfazz/go-apt/example/fazzdb_sample/migration"
	"github.com/payfazz/go-apt/pkg/fazzdb"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	fazzdb.Migrate(config.GetDB(), "cashfazz-example",
		migration.Version1,
		migration.Version2,
	)

	query := fazzdb.QueryDb(config.GetDB(), config.Parameter)
	//fazzdb_sample.InsertAuthor(query)
	//fazzdb_sample.InsertBook(query)
	//
	//fazzdb_sample.UpdateAuthor(query)
	//fazzdb_sample.UpdateBook(query)
	//
	//fazzdb_sample.BulkInsertAuthors(query)
	//fazzdb_sample.BulkInsertBooks(query)

	fazzdb_sample.FirstAuthor(query)
	fazzdb_sample.FirstBook(query)

	//fazzdb_sample.DeleteBook(query)
	//fazzdb_sample.DeleteAuthor(query)

	fazzdb_sample.AllAuthors(query)
	fazzdb_sample.AllBooks(query)

	fazzdb_sample.RawFirst(query)
	fazzdb_sample.RawAll(query)

	fazzdb_sample.RawNamedFirst(query)
	fazzdb_sample.RawNamedAll(query)

	fazzdb_sample.Sum(query)
	fazzdb_sample.Count(query)
	fazzdb_sample.Max(query)
	fazzdb_sample.Min(query)
	fazzdb_sample.Avg(query)
}

/*func RawFirst(query *fazzdb.Query) {
	qry := "SELECT name, age FROM students WHERE age > $1 AND age < $2 LIMIT 1;"
	payload := []interface{}{
		20,
		25,
	}

	r, err := query.RawFirst(&model.StudentCompact{}, qry, payload...)
	if nil != err {
		_ = query.Tx.Rollback()
		panic(err)
	}
	result := r.(model.StudentCompact)

	fmt.Println(result.Name, "-", result.Age)
}

func RawAll(query *fazzdb.Query) {
	qry := "SELECT name, age FROM students WHERE age > $1 AND age < $2;"
	payload := []interface{}{
		20,
		25,
	}

	r, err := query.RawAll(&model.StudentCompact{}, qry, payload...)
	if nil != err {
		_ = query.Tx.Rollback()
		panic(err)
	}
	results := r.([]model.StudentCompact)

	for _, v := range results {
		fmt.Println(v.Name, "-", v.Age)
	}
}

func RawNamedExec(query *fazzdb.Query) {
	payload := map[string]interface{}{}
	payload["name"] = "Jono"
	payload["address"] = "Dadap"
	payload["age"] = 30

	_, err := query.RawNamedExec("INSERT INTO students (name, address, age) VALUES (:name, :address, :age);", payload)
	if nil != err {
		_ = query.Tx.Rollback()
		panic(err)
	}
}

func RawNamedFirst(query *fazzdb.Query) {
	qry := "SELECT name, age FROM students WHERE age > :age LIMIT 1;"
	payload := map[string]interface{}{}
	payload["age"] = 25

	r, err := query.RawNamedFirst(&model.StudentCompact{}, qry, payload)
	if nil != err {
		_ = query.Tx.Rollback()
		panic(err)
	}

	result := r.(model.StudentCompact)

	fmt.Println(result.Name, "-", result.Age)
}

func RawNamedAll(query *fazzdb.Query) {
	qry := "SELECT name, age FROM students WHERE age > :age;"
	payload := map[string]interface{}{}
	payload["age"] = 25

	r, err := query.RawNamedAll(&model.StudentCompact{}, qry, payload)
	if nil != err {
		_ = query.Tx.Rollback()
		panic(err)
	}

	results := r.([]model.StudentCompact)

	for _, v := range results {
		fmt.Println(v.Name, "-", v.Age)
	}
}

func Sum(query *fazzdb.Query, column string) {
	n := model.NewStudent()
	result, err := query.Use(n).
		Sum(column)

	if nil != err {
		_ = query.Tx.Rollback()
		panic(err)
	}

	fmt.Println("SUM: ", *result)
}*/