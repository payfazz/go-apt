package model

import "github.com/payfazz/go-apt/pkg/fazzdb"

type Student struct {
	*fazzdb.Model
	Id      int    `db:"id"`
	Name    string `db:"name"`
	Address string `db:"address"`
	Age     int    `db:"age"`
}

func (s *Student) Get(key string) interface{} {
	return s.Payload()[key]
}

func (s *Student) Payload() map[string]interface{} {
	return s.MapPayload(s)
}

func NewStudent() *Student {
	model := fazzdb.AutoIncrementModel(
		"students",
		[]string{
			"id",
			"name",
			"address",
			"age",
		},
		"id",
		false,
		false,
	)
	return &Student{
		Model: model,
	}
}
