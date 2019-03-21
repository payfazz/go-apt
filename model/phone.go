package model

import "github.com/payfazz/go-apt/pkg/fazzdb"

type Phone struct {
	*fazzdb.Model
	Id  int    `db:"id"`
	Num string `db:"num"`
}

func (s *Phone) Get(key string) interface{} {
	return s.Payload()[key]
}

func (s *Phone) Payload() map[string]interface{} {
	return s.MapPayload(s)
}

func NewPhone() *Phone {
	model := fazzdb.AutoIncrementModel(
		"phones",
		[]string{
			"id",
			"num",
		},
		"id",
		false,
		false,
	)
	return &Phone{
		Model: model,
	}
}