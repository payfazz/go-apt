package seed

import (
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

type AuthorSeed struct {
	fazzdb.Seeder
}

func (b *AuthorSeed) Type() fazzdb.SeedType {
	return fazzdb.SEED_RAW
}

func (a *AuthorSeed) RawQuery() string {
	return `INSERT INTO authors ("name", "country") VALUES ` +
		`('Charles Dickens', 'United States'),` +
		`('JK Rowling', 'United Kingdom'),` +
		`('Fitzgerald', 'German');`
}

func AuthorSeeder() *AuthorSeed {
	return &AuthorSeed{}
}
