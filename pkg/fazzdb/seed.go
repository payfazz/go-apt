package fazzdb

type SeedType string

const (
	SEED_RAW  SeedType = "RAW"
	SEED_BULK SeedType = "BULK"
)

type SeederInterface interface {
	Type() SeedType
	Model() ModelInterface
	BulkModels() []ModelInterface
	RawQuery() string
}

type Seeder struct {
	Seeds []map[string]interface{}
}

func (s *Seeder) Type() SeedType {
	panic("Please override Type() method")
}

func (s *Seeder) Model() ModelInterface {
	panic("Please override Model() method")
}

func (s *Seeder) BulkModels() []ModelInterface {
	panic("Please override BulkModels() method")
}

func (s *Seeder) RawQuery() string {
	panic("Please override RawQuery() method")
}

func Seed(query *Query, seeders ...SeederInterface) {
	var err error
	for _, seeder := range seeders {
		if SEED_BULK == seeder.Type() {
			_, err = query.Use(seeder.Model()).BulkInsert(seeder.BulkModels())
		} else if SEED_RAW == seeder.Type() {
			_, err = query.RawExec(seeder.RawQuery())
		}

		if nil != err {
			_ = query.Tx.Rollback()
			panic(err)
		}
	}
}
