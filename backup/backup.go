package backup

import (
	"github.com/jmoiron/sqlx"
	"reflect"
)

func getStudent(db *sqlx.DB, result interface{}) {
	row := db.QueryRowx("SELECT * FROM students LIMIT 1")
	_ = row.StructScan(result)
}

func getStudentImmutable(db *sqlx.DB, sample interface{}) interface{} {
	var temp = reflect.New(reflect.TypeOf(sample)).Interface()
	row := db.QueryRowx("SELECT * FROM students LIMIT 1")
	_ = row.StructScan(temp)

	return temp
}

func getStudentsImmutable(db *sqlx.DB, sample interface{}) interface{} {
	var results = reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(sample)), 0, 0)

	rows, _ := db.Queryx("SELECT * FROM students")
	for rows.Next() {
		var temp = reflect.New(reflect.TypeOf(sample)).Interface()
		_ = rows.StructScan(temp)
		results = reflect.Append(results, reflect.ValueOf(temp).Elem())
	}

	return results.Interface()
}
