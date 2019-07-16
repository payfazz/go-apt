package main

import (
	"math/rand"
	"time"

	_ "github.com/lib/pq"
	"github.com/payfazz/go-apt/config"
	"github.com/payfazz/go-apt/example/db/fazzdb_sample"
	"github.com/payfazz/go-apt/example/db/fazzdb_sample/migration"
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	fazzdb.Verbose()
	fazzdb.Migrate(config.GetDB(), "cashfazz-example", true, true,
		migration.Version1,
		migration.Version2,
	)

	query := fazzdb.QueryDb(config.GetDB(), config.Parameter)
	fazzdb_sample.InsertAuthor(query)
	fazzdb_sample.InsertBook(query)

	fazzdb_sample.UpdateAuthor(query)
	fazzdb_sample.UpdateBook(query)

	fazzdb_sample.BulkInsertAuthors(query)
	fazzdb_sample.BulkInsertBooks(query)

	fazzdb_sample.FirstAuthor(query)
	fazzdb_sample.FirstBook(query)

	fazzdb_sample.DeleteBook(query)
	fazzdb_sample.DeleteAuthor(query)

	fazzdb_sample.AllAuthors(query)
	fazzdb_sample.AllBooks(query)
	fazzdb_sample.AllBooksSliceConditions(query)

	fazzdb_sample.RawFirst(query)
	fazzdb_sample.RawAll(query)

	fazzdb_sample.RawNamedFirst(query)
	fazzdb_sample.RawNamedAll(query)

	fazzdb_sample.Sum(query)
	fazzdb_sample.Count(query)
	fazzdb_sample.Max(query)
	fazzdb_sample.Min(query)
	fazzdb_sample.Avg(query)
}
