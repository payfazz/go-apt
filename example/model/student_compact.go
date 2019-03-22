package model

type StudentCompact struct {
	Name string `db:"name"`
	Age  int    `db:"age"`
}
