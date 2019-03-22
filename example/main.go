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

	//arg := map[string]interface{}{
	//	"published": true,
	//	"authors": []interface{}{8, 19, 32, 44},
	//}
	//log.Println(arg)
	//a, args, _ := sqlx.Named("SELECT * FROM articles WHERE published=:published AND author_id IN (:authors)", arg)
	//log.Println(a)
	//log.Println(args)
	//a, args, _ = sqlx.In(a, args...)
	//log.Println(a)
	//log.Println(args)
	//a = db.Rebind(a)
	//log.Println(a)

	//n := model.NewUid()
	//n.Data = 10
	//_, err := query.Use(n).
	//	Insert()
	//if nil != err {
	//	log.Println(err)
	//}
	//
	//fmt.Println(n.Id)

	// Select Many
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

	//p := model.NewPhone()
	//_, err := query.Use(p).
	//	GetAll()
	//
	//if nil != err {
	//	panic(err)
	//}

	//
	//phones := results2.([]model.Phone)
	//for _, p2 := range phones {
	//	fmt.Println(p2.Model)
	//}

	// Select One
	//n := model.NewStudent()
	//result, err := query.Use(n).
	//	First()
	//
	//if nil != err {
	//	panic(err)
	//}
	//student := result.(model.Student)
	//fmt.Printf("%d - %s - %s - %d\n", student.Id, student.Name, student.Address, student.Age)
	//
	//student.Name = "Hi"
	//_, err = query.Use(&student).Update()
	//if nil != err {
	//	panic(err)
	//}

	// Delete
	//_, _ = query.Use(&student).
	//	Delete()

	_ = tx.Commit()
}
