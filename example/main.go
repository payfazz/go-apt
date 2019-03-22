package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/payfazz/go-apt/example/config"
	"github.com/payfazz/go-apt/example/model"
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

func main() {
	conn := "host=localhost port=5432 user=postgres password=cashfazz dbname=qb sslmode=disable"

	db, _ := sqlx.Connect("postgres", conn)
	tx, _ := db.Beginx()
	query := fazzdb.QueryTx(tx, config.Db)

	//student := SelectOne(query)
	//Update(query, student)

	Aggregate(query, fazzdb.AG_COUNT, "age")

	_ = tx.Commit()
}

func SelectOne(query *fazzdb.Query) *model.Student {
	// Select One
	n := model.NewStudent()
	result, err := query.Use(n).
		First()

	if nil != err {
		panic(err)
	}
	student := result.(model.Student)
	fmt.Printf("%d - %s - %s - %d\n", student.Id, student.Name, student.Address, student.Age)

	return &student
}

func Update(query *fazzdb.Query, student *model.Student) {
	student.Name = "Hi"
	_, err := query.Use(student).Update()
	if nil != err {
		panic(err)
	}
}

func Delete(query *fazzdb.Query, student *model.Student) {
	_, _ = query.Use(student).
		Delete()
}

func SelectMany(query *fazzdb.Query) {
	n := model.NewStudent()
	results, err := query.Use(n).
		GroupWhere(func(query *fazzdb.Query) *fazzdb.Query {
			return query.WhereOp("name", fazzdb.OP_LIKE, "%i%").
				WhereOp("age", fazzdb.OP_IN, []interface{}{22, 23, 24, 25})
		}).
		WhereOp("address", fazzdb.OP_LIKE, "%i%").
		OrderBy("age", fazzdb.DIR_DESC).
		GetAll()

	if nil != err {
		panic(err)
	}

	students := results.([]model.Student)
	for _, s := range students {
		fmt.Printf("%d - %s - %s - %d\n", s.Id, s.Name, s.Address, s.Age)
	}
}

func Aggregate(query *fazzdb.Query, aggregate fazzdb.Aggregate, column string) {
	n := model.NewStudent()
	result, err := query.Use(n).
		Aggregate(aggregate, column)

	if nil != err {
		panic(err)
	}

	fmt.Println(aggregate, ": ", *result)
}