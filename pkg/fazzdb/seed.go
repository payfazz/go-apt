package fazzdb

import (
	"time"

	"github.com/gofrs/uuid"
)

type SeedType string
type PKType string

const (
	// SEED_RAW seed using raw query
	SEED_RAW SeedType = "RAW"
	// SEED_BULK seed using bulk model
	SEED_BULK SeedType = "BULK"
	// SEED_OBJECT seed using map[string]interface{} object
	SEED_OBJECT SeedType = "OBJECT"
)

const (
	// PK_UUID use if primary key type is uuid
	PK_UUID PKType = "UUID"
	// PK_UUID use if primary key type is serial / big serial / auto increment
	PK_SERIAL PKType = "SERIAL"
	// PK_PLAIN use if primary key type is not listed
	PK_PLAIN PKType = "PLAIN"
)

// SeederInterface is an interface used for creating seeder
type SeederInterface interface {
	Type() SeedType
	Table() TableMeta
	Values() []map[string]interface{}
	Model() ModelInterface
	BulkModels() []ModelInterface
	RawQuery() string
}

// TableMeta is a struct to store table meta data for Seeder
type TableMeta struct {
	Name       string
	PrimaryKey string
	Type       PKType
	CreatedAt  bool
	Columns    []string
}

// NewTableMeta is a constructor to create table meta
func NewTableMeta(name string, primaryKey string, pkType PKType, createdAt bool, columns []string) TableMeta {
	pkExist := false
	createdAtExist := false
	for _, v := range columns {
		if v == primaryKey {
			pkExist = true
		}
		if v == CREATED_AT {
			createdAtExist = true
		}
	}

	if !pkExist && PK_UUID == pkType {
		columns = append(columns, primaryKey)
	}

	if !createdAtExist {
		columns = append(columns, CREATED_AT)
	}

	return TableMeta{
		Name:       name,
		PrimaryKey: primaryKey,
		Type:       pkType,
		CreatedAt:  createdAt,
		Columns:    columns,
	}
}

// Seeder is a base struct for seeder with base implementation of Seeder data
type Seeder struct {
	Seeds []map[string]interface{}
}

// Type is a function that will return SeedType (SEED_RAW / SEED_OBJECT / SEED_BULK)
func (s *Seeder) Type() SeedType {
	panic("Please override Type() method")
}

// Table is a function that will return table name when using SEED_OBJECT
func (s *Seeder) Table() TableMeta {
	panic("Please override Table() method")
}

// Values is a function that will return slice of map[string]interface{} with seed data when using SEED_OBJECT
func (s *Seeder) Values() []map[string]interface{} {
	panic("Please override Values() method")
}

// Model is a function that will return model instance used for seeder when using SEED_BULK
func (s *Seeder) Model() ModelInterface {
	panic("Please override Model() method")
}

// BulkModels is a function that will return slices of model with seed data when using SEED_BULK
func (s *Seeder) BulkModels() []ModelInterface {
	panic("Please override BulkModels() method")
}

// RawQuery is a function that will contain raw query used for seeder when using SEED_RAW
func (s *Seeder) RawQuery() string {
	panic("Please override RawQuery() method")
}

// generateEmpty is a function that will generate uuid primary key if not specified
func generateEmpty(meta TableMeta, values []map[string]interface{}) []map[string]interface{} {
	if PK_UUID != meta.Type {
		return values
	}

	for i, v := range values {
		if _, ok := v[meta.PrimaryKey]; !ok || nil == v[meta.PrimaryKey] {
			v4, _ := uuid.NewV4()
			values[i][meta.PrimaryKey] = v4.String()
		}
		if _, ok := v[CREATED_AT]; !ok || nil == v[CREATED_AT] {
			now := time.Now().UTC()
			values[i][CREATED_AT] = now.Format(time.RFC3339)
		}
	}

	return values
}

// Seed is a function for entry point for running seeder
func Seed(query *Query, seeders ...SeederInterface) {
	var err error
	for _, seeder := range seeders {
		if SEED_OBJECT == seeder.Type() {
			builder := NewBuilder()
			meta := seeder.Table()
			values := generateEmpty(meta, seeder.Values())
			seedQuery := builder.BuildSeeder(meta.Name, meta.Columns, values)

			_, err = query.RawExec(seedQuery)
		} else if SEED_BULK == seeder.Type() {
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
