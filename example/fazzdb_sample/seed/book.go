package seed

import (
	"github.com/payfazz/go-apt/example/fazzdb_sample/model"
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

type BookSeed struct {
	fazzdb.Seeder
}

func (b *BookSeed) Type() fazzdb.SeedType {
	return fazzdb.SEED_BULK
}

func (b *BookSeed) Model() fazzdb.ModelInterface {
	return model.BookModel()
}

func (b *BookSeed) BulkModels() []fazzdb.ModelInterface {
	var bulks []fazzdb.ModelInterface
	for _, v := range b.Seeds {
		book := model.BookModel()
		book.Title = v["title"].(string)
		book.Year = v["year"].(int)
		book.Stock = v["stock"].(int)
		book.Status = v["status"].(string)
		book.AuthorId = v["authorId"].(int)

		bulks = append(bulks, book)
	}
	return bulks
}

func BookSeeder() *BookSeed {
	return &BookSeed{
		Seeder: fazzdb.Seeder{
			Seeds: []map[string]interface{}{
				{"title": "Harry Motter", "year": 2000, "stock": 20, "status": model.BOOK_AVAILABLE, "authorId": 1},
				{"title": "Subtle Art", "year": 2002, "stock": 15, "status": model.BOOK_AVAILABLE, "authorId": 2},
				{"title": "Dumbo", "year": 2003, "stock": 52, "status": model.BOOK_AVAILABLE, "authorId": 2},
				{"title": "Aladin", "year": 2005, "stock": 0, "status": model.BOOK_BORROWED, "authorId": 3},
				{"title": "Iron men", "year": 1990, "stock": 10, "status": model.BOOK_AVAILABLE, "authorId": 1},
				{"title": "Blackjack", "year": 1967, "stock": 15, "status": model.BOOK_AVAILABLE, "authorId": 1},
				{"title": "21", "year": 1972, "stock": 18, "status": model.BOOK_AVAILABLE, "authorId": 2},
				{"title": "Eagle Eye", "year": 1978, "stock": 12, "status": model.BOOK_AVAILABLE, "authorId": 2},
				{"title": "Hurk", "year": 1982, "stock": 22, "status": model.BOOK_AVAILABLE, "authorId": 3},
				{"title": "Viking", "year": 1989, "stock": 25, "status": model.BOOK_AVAILABLE, "authorId": 3},
				{"title": "Tor", "year": 1993, "stock": 8, "status": model.BOOK_AVAILABLE, "authorId": 1},
				{"title": "Roki", "year": 1996, "stock": 13, "status": model.BOOK_AVAILABLE, "authorId": 2},
			},
		},
	}
}
