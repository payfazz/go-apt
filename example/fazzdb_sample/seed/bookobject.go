package seed

import (
	"github.com/payfazz/go-apt/example/fazzdb_sample/model"
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

type BookObjectSeed struct {
	fazzdb.Seeder
}

func (b *BookObjectSeed) Type() fazzdb.SeedType {
	return fazzdb.SEED_OBJECT
}

func (b *BookObjectSeed) Table() fazzdb.TableMeta {
	return fazzdb.NewTableMeta(
		"books",
		"id",
		fazzdb.PK_UUID,
		[]string{
			"title",
			"year",
			"stock",
			"status",
			"authorId",
		},
	)
}

func (b *BookObjectSeed) Values() []map[string]interface{} {
	return []map[string]interface{}{
		{"title": "Eagle Eye", "year": 1978, "stock": 12, "status": model.BOOK_AVAILABLE, "authorId": 2},
		{"title": "Hurk", "year": 1982, "stock": 22, "status": model.BOOK_AVAILABLE, "authorId": 3},
		{"title": "Viking", "year": 1989, "stock": 25, "status": model.BOOK_AVAILABLE, "authorId": 3},
		{"title": "Tor", "year": 1993, "stock": 8, "status": model.BOOK_AVAILABLE, "authorId": 1},
		{"title": "Roki", "year": 1996, "stock": 13, "status": model.BOOK_AVAILABLE, "authorId": 2},
	}
}

func BookObjectSeeder() *BookObjectSeed {
	return &BookObjectSeed{}
}
