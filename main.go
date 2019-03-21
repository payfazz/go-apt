package main

import (
	"db/fazzdb"
	"db/model"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	conn := "host=localhost port=5432 user=postgres password=cashfazz dbname=qb sslmode=disable"

	db, _ := sqlx.Connect("postgres", conn)
	tx, _ := db.Beginx()

	//results, err := fazzdb.NewQueryTx(tx, model.NewStudent()).
	//	OrderBy("age", fazzorder.DIR_DESC).
	//	GetAll()
	//if nil != err {
	//	panic(err)
	//}
	//students := results.([]model.Student)
	//
	//for _, s := range students {
	//	fmt.Printf("%d - %s - %s - %d\n", s.Id, s.Name, s.Address, s.Age)
	//}
	//
	//fmt.Println("------")
	//
	//result, err := fazzdb.NewQuery(db, model.NewStudent()).
	//	Where("name", "asd").
	//	First()
	//if nil == result {
	//	log.Println("No data")
	//	return
	//}
	//if nil != err {
	//	panic(err)
	//}
	//student := result.(model.Student)
	//
	//fmt.Printf("%d - %s - %s - %d\n", student.Id, student.Name, student.Address, student.Age)

	//student := model.NewStudent()
	//student.Name = "Test2"
	//student.Address = "Test 2Address"
	//student.Age = 17
	//_, _ = fazzdb.NewQueryTx(tx, student).Insert()

	//student := model.NewStudent()
	//student.Payload()

	_, err := fazzdb.NewQueryTx(tx, model.NewStudent()).
		GetAll()
	if nil != err {
		panic(err)
	}

//	student.Name = "AbcAndi"
//	_, _ = fazzdb.NewQueryTx(tx, &student).Update()

	_ = tx.Commit()
}


