package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/payfazz/go-apt/config"
	"github.com/payfazz/go-apt/example/model"
	"github.com/payfazz/go-apt/pkg/fazzdb"
	"log"
	"math/rand"
)

func main() {
	conn := "host=localhost port=5432 user=postgres password=cashfazz dbname=qb sslmode=disable"

	db, _ := sqlx.Connect("postgres", conn)
	tx, _ := db.Beginx()
	//query := fazzdb.QueryTx(tx, config.Db)
	query := fazzdb.QueryDb(db, config.Db)

	//Insert(query)
	//RawFirst(query)
	//RawAll(query)
	//student := SelectOne(query)
	//Delete(query, student)
	//Update(query, student)

	SelectAll(query)
	//SelectOne(query)

	_ = tx.Commit()
}

func BulkInsert(query *fazzdb.Query) {
	students := make([]*model.Student, 0)
	for i := 0; i < 20; i++ {
		student := model.NewStudent()
		student.Name = fmt.Sprintf("Bulk%d", i)
		student.Address = fmt.Sprintf("Address %d", i)
		student.Age = rand.Intn(20) + 12

		students = append(students, student)
	}

	_, err := query.Use(model.NewStudent()).BulkInsert(students)
	if nil != err {
		panic(err)
	}
}

func RawFirst(query *fazzdb.Query) {
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

func SelectOne(query *fazzdb.Query) *model.Student {
	// Select One
	n := model.NewStudent()
	result, err := query.Use(n).
		First()

	if nil != err {
		_ = query.Tx.Rollback()
		panic(err)
	}

	student := result.(model.Student)
	fmt.Printf("%d - %s - %s - %d - %s\n", student.Id, student.Name, student.Address, student.Age, student.CreatedAt)

	return &student
}

func Insert(query *fazzdb.Query) {
	student := model.NewStudent()
	student.Name = "sby"
	student.Address = "Solo"
	student.Age = 20

	id, err := query.Use(student).Insert()
	if nil != err {
		panic(err)
	}
	log.Println("Inserted id:", *id)
}

func Update(query *fazzdb.Query, student *model.Student) {
	student.Name = "Hi123"
	_, err := query.Use(student).Update()
	if nil != err {
		_ = query.Tx.Rollback()
		panic(err)
	}
}

func Delete(query *fazzdb.Query, student *model.Student) {
	_, err := query.Use(student).
		Delete()
	if nil != err {
		_ = query.Tx.Rollback()
		panic(err)
	}
}

func SelectAll(query *fazzdb.Query) {
	n := model.NewStudent()
	results, err := query.Use(n).
		All()

	if nil != err {
		_ = query.Tx.Rollback()
		panic(err)
	}

	students := results.([]model.Student)
	for _, s := range students {
		fmt.Printf("%d - %s - %s - %d - %s\n", s.Id, s.Name, s.Address, s.Age, s.CreatedAt)
		//fmt.Println(reflect.TypeOf(s.CreatedAt))
	}
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
		All()

	if nil != err {
		_ = query.Tx.Rollback()
		panic(err)
	}

	students := results.([]model.Student)
	for _, s := range students {
		fmt.Printf("%d - %s - %s - %d - %s - %s - %s\n", s.Id, s.Name, s.Address, s.Age, s.CreatedAt, s.UpdatedAt, s.DeletedAt)
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
}