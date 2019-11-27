package fazzdb

import (
	"reflect"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/payfazz/go-apt/pkg/fazzcommon/formatter"
)

const (
	MODEL_STRUCT_NAME    = "Model"
	TIME_PTR_STRUCT_NAME = "*time.Time"
	TIME_STRUCT_NAME     = "time.Time"
)

// ModelInterface is an interface that will be used to get model information and used in various task by Query instance
type ModelInterface interface {
	// GetModel is a function that will return pointer to Model instance
	GetModel() Model
	// GetTable is a function that will return the table name of the Model instance
	GetTable() string
	// GetCreatedAt is a function that will return created_at value
	GetCreatedAt() *time.Time
	// GetUpdatedAt is a function that will return updated_at value
	GetUpdatedAt() *time.Time
	// GetDeletedAt is a function that will return deleted_at value
	GetDeletedAt() *time.Time
	// GetColumns is a function that will return the slice of columns of the Model instance
	GetColumns() []Column
	// GetPK is a function that will return the primary key field name of the Model instance
	GetPK() string
	// Get is a function that MUST be overridden by all model, if it's not overridden it will panic.
	Get(key string) interface{}
	// GeneratePK is a function that MUST be overridden by UuidModel, if it's not overridden it will panic.
	GeneratePK()
	// GenerateId is a function that will generate uuid for primary key if the model created using UuidModel constructor
	GenerateId(v ModelInterface)
	// ColumnCount is a function that will return the length of columns of the Model instance
	ColumnCount() int
	// IsTimestamps is a function that will return true if the Model instance using created_at and updated_at field
	IsTimestamps() bool
	// IsSoftDelete is a function that will return true if the Model instance is using deleted_at field
	IsSoftDelete() bool
	// IsUuid is a function that will return true if the Model instance is created using UuidModel constructor
	IsUuid() bool
	// IsAutoIncrement is a function that will return true if the Model instance is created using AutoIncrementModel constructor
	IsAutoIncrement() bool
	// Payload is a function that MUST be overridden by all model, if it's not overridden it will panic.
	Payload() map[string]interface{}
	// MapPayload is a function that will map all column value as a map[string]interface{} with lowered first character as key
	MapPayload(v interface{}) map[string]interface{}

	// created is a function that will set created_at field with current time, used when inserting model with timestamp
	created()
	// updated is a function that will set updated_at field with current time, used when updating model with timestamp
	updated()
	// deleted is a function that will set deleted_at field with current time, used when soft deleting model
	deleted()
	// recovered is a function that will set deleted_at field with nil, used when recovering soft deleted model
	recovered()
	// applyTimestamps is a function that will copy CreatedAt/UpdatedAt/DeletedAt value into respective timestamp type
	applyTimestamps()
}

// UuidModel is a constructor that is used to initialize a new model with uuid as primary key
func UuidModel(table string, columns []Column, primaryKey string, timestamps bool, softDelete bool) Model {
	return newModel(table, columns, primaryKey, timestamps, softDelete, true, false)
}

// PlainModel is a constructor that is used to initialize a new model with primary key that is neither
// uuid or autoincrement
func PlainModel(table string, columns []Column, primaryKey string, timestamps bool, softDelete bool) Model {
	return newModel(table, columns, primaryKey, timestamps, softDelete, false, false)
}

// AutoIncrementModel is a constructor that is used to initialize a new model with autoincrement primary key
func AutoIncrementModel(table string, columns []Column, primaryKey string, timestamps bool, softDelete bool) Model {
	return newModel(table, columns, primaryKey, timestamps, softDelete, false, true)
}

// UuidSnakeModel is a constructor that is used to initialize a new model with uuid as primary key
func UuidSnakeModel(table string, columns []Column, primaryKey string, timestamps bool, softDelete bool) Model {
	return newSnakeModel(table, columns, primaryKey, timestamps, softDelete, true, false)
}

// PlainSnakeModel is a constructor that is used to initialize a new model with primary key that is neither
// uuid or autoincrement
func PlainSnakeModel(table string, columns []Column, primaryKey string, timestamps bool, softDelete bool) Model {
	return newSnakeModel(table, columns, primaryKey, timestamps, softDelete, false, false)
}

// AutoIncrementSnakeModel is a constructor that is used to initialize a new model with autoincrement primary key
func AutoIncrementSnakeModel(table string, columns []Column, primaryKey string, timestamps bool, softDelete bool) Model {
	return newSnakeModel(table, columns, primaryKey, timestamps, softDelete, false, true)
}

// newModel is a base constructor that will return Model instance that will be used by
// uuid / plain / autoincrement model
func newModel(
	table string,
	columns []Column,
	primaryKey string,
	timestamps bool,
	softDelete bool,
	isUuid bool,
	isAutoIncrement bool,
) Model {
	model := Model{
		Table:         table,
		Columns:       columns,
		PrimaryKey:    primaryKey,
		Uuid:          isUuid,
		AutoIncrement: isAutoIncrement,
		Timestamps:    timestamps,
		SoftDelete:    softDelete,
		TimestampType: TS_CAMEL,
	}

	model.handleTimestamp()
	model.handleSoftDelete()

	return model
}

// newSnakeModel is a base constructor that will return Model instance that will be used by
// uuid / plain / autoincrement model
func newSnakeModel(
	table string,
	columns []Column,
	primaryKey string,
	timestamps bool,
	softDelete bool,
	isUuid bool,
	isAutoIncrement bool,
) Model {
	model := Model{
		Table:         table,
		Columns:       columns,
		PrimaryKey:    primaryKey,
		Uuid:          isUuid,
		AutoIncrement: isAutoIncrement,
		Timestamps:    timestamps,
		SoftDelete:    softDelete,
		TimestampType: TS_SNAKE,
	}

	model.handleTimestamp()
	model.handleSoftDelete()

	return model
}

// Model is a struct that defines the base requirement for a model that will be made, it includes Timestamps and
// SoftDelete field that will be available if it's needed and ignored when not needed
type Model struct {
	Table         string     `json:"-"`
	Columns       []Column   `json:"-"`
	PrimaryKey    string     `json:"-"`
	Uuid          bool       `json:"-"`
	AutoIncrement bool       `json:"-"`
	Timestamps    bool       `json:"-"`
	SoftDelete    bool       `json:"-"`
	TimestampType string     `json:"-"`
	CreatedAt     *time.Time `json:"-" db:"created_at"`
	UpdatedAt     *time.Time `json:"-" db:"updated_at"`
	DeletedAt     *time.Time `json:"-" db:"deleted_at"`
	CamelTimestamp
	SnakeTimestamp
}

// GeneratePK is a function that MUST be overridden by UuidModel, if it's not overridden it will panic.
// The overriding function only need to call GenerateId(v interface{}) function with its own struct as
// the parameter:
//
//  func (s *Student) GeneratePK() {
//      s.GenerateId(s)
//  }
func (m *Model) GeneratePK() {
	if m.IsUuid() {
		panic("Please override GeneratePK() method in your model")
	}
}

// Get is a function that MUST be overridden by all model, if it's not overridden it will panic.
// The overriding function only need to call the array of map[string]interface{} provided by Payload() function:
//
//  func (s *Uid) Get(key string) interface{} {
//      return s.Payload()[key]
//  }
func (m *Model) Get(key string) interface{} {
	panic("Please override Get(key string) method in your model")
}

// Payload is a function that MUST be overridden by all model, if it's not overridden it will panic.
// The overriding function only need to call MapPayload(v interface{}) function with its own struct as
// the parameter:
//
//  func (s *Uid) Payload() map[string]interface{} {
//      return s.MapPayload(s)
//  }
func (m *Model) Payload() map[string]interface{} {
	panic("Please override Payload() method in your model")
}

// GetModel is a function that will return pointer to Model instance
func (m *Model) GetModel() Model {
	return *m
}

// GetTable is a function that will return the table name of the Model instance
func (m *Model) GetTable() string {
	return m.Table
}

// GetCreatedAt is a function that will return created_at value
func (m *Model) GetCreatedAt() *time.Time {
	return m.CreatedAt
}

// GetUpdatedAt is a function that will return updated_at value
func (m *Model) GetUpdatedAt() *time.Time {
	return m.UpdatedAt
}

// GetDeletedAt is a function that will return deleted_at value
func (m *Model) GetDeletedAt() *time.Time {
	return m.DeletedAt
}

// GetColumns is a function that will return the slice of columns of the Model instance
func (m *Model) GetColumns() []Column {
	return m.Columns
}

// GetPK is a function that will return the primary key field name of the Model instance
func (m *Model) GetPK() string {
	return m.PrimaryKey
}

// GenerateId is a function that will generate uuid for primary key if the model created
// using UuidModel constructor
func (m *Model) GenerateId(v ModelInterface) {
	if !m.Uuid || "" != v.Get(v.GetPK()) {
		return
	}

	pkField := strings.Title(m.GetPK())

	v4, _ := uuid.NewV4()
	id := v4.String()

	reflect.ValueOf(v).Elem().FieldByName(pkField).Set(reflect.ValueOf(id))
}

// ColumnCount is a function that will return the length of columns of the Model instance
func (m *Model) ColumnCount() int {
	return len(m.Columns)
}

// IsTimestamps is a function that will return true if the Model instance using created_at and updated_at field
func (m *Model) IsTimestamps() bool {
	return m.Timestamps
}

// IsSoftDelete is a function that will return true if the Model instance is using deleted_at field
func (m *Model) IsSoftDelete() bool {
	return m.SoftDelete
}

// IsUuid is a function that will return true if the Model instance is created using UuidModel constructor
func (m *Model) IsUuid() bool {
	return m.Uuid
}

// IsAutoIncrement is a function that will return true if the Model instance is created using AutoIncrementModel
// constructor
func (m *Model) IsAutoIncrement() bool {
	return m.AutoIncrement
}

// MapPayload is a function that will map all column value as a map[string]interface{} with lowered first character as key
func (m *Model) MapPayload(v interface{}) map[string]interface{} {
	var results = make(map[string]interface{})
	classValue := reflect.ValueOf(v)
	classType := reflect.TypeOf(v)

	if reflect.Ptr == classValue.Kind() {
		classValue = classValue.Elem()
		classType = classType.Elem()
	}

	for i := 0; i < classValue.Type().NumField(); i++ {
		cValue := classValue.Field(i)
		cType := classType.Field(i)
		tag := formatter.ToLowerFirst(cType.Tag.Get("db"))

		if reflect.Ptr == cValue.Kind() && TIME_PTR_STRUCT_NAME != cValue.Type().String() {
			cValue = cValue.Elem()
		}

		if reflect.Invalid == cValue.Kind() && "" == tag {
			continue
		}

		if reflect.Invalid == cValue.Kind() {
			results[tag] = nil
		} else if TIME_PTR_STRUCT_NAME == cValue.Type().String() || TIME_STRUCT_NAME == cValue.Type().String() {
			results[tag] = cValue.Interface()
		} else if MODEL_STRUCT_NAME == cValue.Type().Name() {
			model := cValue.Interface().(Model)
			if model.IsTimestamps() {
				results[CREATED_AT] = model.CreatedAt
				results[UPDATED_AT] = model.UpdatedAt
			}
			if model.IsSoftDelete() {
				results[DELETED_AT] = model.DeletedAt
			}
		} else if reflect.Struct == cValue.Kind() {
			mapResults := m.MapPayload(cValue.Interface())
			for mapIndex, mapValue := range mapResults {
				results[mapIndex] = mapValue
			}
		} else if "" != tag {
			results[tag] = cValue.Interface()
		}
	}
	return results
}

// created is a function that will set created_at field with current time, used when inserting model with timestamp
func (m *Model) created() {
	if m.IsTimestamps() && nil == m.CreatedAt {
		now := time.Now()
		m.CreatedAt = &now
		m.applyTimestamps()
	}
}

// updated is a function that will set updated_at field with current time, used when updating model with timestamp
func (m *Model) updated() {
	if m.IsTimestamps() {
		now := time.Now()
		m.UpdatedAt = &now
		m.applyTimestamps()
	}
}

// deleted is a function that will set deleted_at field with current time, used when soft deleting model
func (m *Model) deleted() {
	if m.IsSoftDelete() {
		now := time.Now()
		m.DeletedAt = &now
		m.applyTimestamps()
	}
}

// recovered is a function that will set deleted_at field with nil, used when recovering soft deleted model
func (m *Model) recovered() {
	if m.IsSoftDelete() {
		m.DeletedAt = nil
		m.applyTimestamps()
	}
}

// handleTimestamp is a function that will automatically append created_at and updated_at to
// Columns attribute in Model instance
func (m *Model) handleTimestamp() {
	if m.IsTimestamps() {
		m.Columns = append(m.Columns, Col(CREATED_AT))
		m.Columns = append(m.Columns, Col(UPDATED_AT))
	}
}

// handleSoftDelete is a function that will automatically append deleted_at to
// Columns attribute in Model instance
func (m *Model) handleSoftDelete() {
	if m.IsSoftDelete() {
		m.Columns = append(m.Columns, Col(DELETED_AT))
	}
}

func (m *Model) applyTimestamps() {
	switch m.TimestampType {
	case TS_CAMEL:
		m.CamelTimestamp.SetTimestamps(m)
	case TS_SNAKE:
		m.SnakeTimestamp.SetTimestamps(m)
	}
}

// CamelTimestamp is a struct that is used to translate timestamps into camelCase json key
type CamelTimestamp struct {
	CamelCreatedAt *time.Time `json:"createdAt,omitempty"`
	CamelUpdatedAt *time.Time `json:"updatedAt,omitempty"`
	CamelDeletedAt *time.Time `json:"deletedAt,omitempty"`
}

// SetTimestamps is a function to set CamelTimestamp value from model timestamp
func (ts *CamelTimestamp) SetTimestamps(m *Model) {
	ts.CamelCreatedAt = m.CreatedAt
	ts.CamelUpdatedAt = m.UpdatedAt
	ts.CamelDeletedAt = m.DeletedAt
}

// SnakeTimestamp is a struct that is used to translate timestamps into snakeCase json key
type SnakeTimestamp struct {
	SnakeCreatedAt *time.Time `json:"created_at,omitempty"`
	SnakeUpdatedAt *time.Time `json:"updated_at,omitempty"`
	SnakeDeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// SetTimestamps is a function to set SnakeTimestamp value from model timestamp
func (ts *SnakeTimestamp) SetTimestamps(m *Model) {
	ts.SnakeCreatedAt = m.CreatedAt
	ts.SnakeUpdatedAt = m.UpdatedAt
	ts.SnakeDeletedAt = m.DeletedAt
}
