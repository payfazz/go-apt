package model

import (
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

// Account is struct for account model
type Account struct {
	fazzdb.Model `json:"-"`
	Id           string `json:"id" db:"id"`
	Version      int    `json:"-" db:"version"`
	Name         string `json:"name" db:"name"`
	Balance      int64  `json:"balance" db:"balance"`
}

// GeneratePK used to generate PK for fazzcard
func (m *Account) GeneratePK() {
	// Don't generate primary key, us
}

// Get used to get the data
func (m *Account) Get(key string) interface{} {
	return m.Payload()[key]
}

// Payload used to get the payload
func (m *Account) Payload() map[string]interface{} {
	return m.MapPayload(m)
}

// AccountModel is a constructor for account model
func AccountModel() *Account {
	return &Account{
		Model: fazzdb.UuidModel(
			"account_read",
			[]fazzdb.Column{
				fazzdb.Col("id"),
				fazzdb.Col("version"),
				fazzdb.Col("name"),
				fazzdb.Col("balance"),
			},
			"id",
			false,
			false,
		),
	}
}
