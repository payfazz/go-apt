package model

import "github.com/payfazz/go-apt/pkg/fazzdb"

type Uid struct {
	*fazzdb.Model
	Id   string `db:"id"`
	Data int    `db:"data"`
}

func (s *Uid) Get(key string) interface{} {
	return s.Payload()[key]
}

func (s *Uid) GeneratePK() {
	s.GenerateId(s)
}

func (s *Uid) Payload() map[string]interface{} {
	return s.MapPayload(s)
}

func NewUid() *Uid {
	model := fazzdb.UuidModel(
		"uids",
		[]string{
			"id",
			"data",
		},
		"id",
		false,
		false,
	)
	return &Uid{
		Model: model,
	}
}
