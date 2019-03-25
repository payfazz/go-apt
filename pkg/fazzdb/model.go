package fazzdb

import (
	"github.com/satori/go.uuid"
	"reflect"
	"strings"
	"time"
	"unicode"
)

type ModelInterface interface {
	GetModel() *Model
	GetTable() string
	GetColumns() []string
	GetPK() string
	Get(key string) interface{}
	GeneratePK()
	GenerateId(v interface{})
	ColumnCount() int
	IsTimestamps() bool
	IsSoftDelete() bool
	IsUuid() bool
	IsAutoIncrement() bool
	Payload() map[string]interface{}
	MapPayload(v interface{}) map[string]interface{}
	Created()
	Updated()
	Deleted()
	Recovered()
}

func UuidModel(table string, columns []string, primaryKey string, timestamps bool, softDelete bool) *Model {
	model := &Model{
		Table:         table,
		Columns:       columns,
		PrimaryKey:    primaryKey,
		Uuid:          true,
		AutoIncrement: false,
		Timestamps:    timestamps,
		SoftDelete:    softDelete,
	}

	model.handleTimestamp()
	model.handleSoftDelete()

	return model
}

func PlainModel(table string, columns []string, primaryKey string, timestamps bool, softDelete bool) *Model {
	model := &Model{
		Table:         table,
		Columns:       columns,
		PrimaryKey:    primaryKey,
		Uuid:          false,
		AutoIncrement: false,
		Timestamps:    timestamps,
		SoftDelete:    softDelete,
	}

	model.handleTimestamp()
	model.handleSoftDelete()

	return model
}

func AutoIncrementModel(table string, columns []string, primaryKey string, timestamps bool, softDelete bool) *Model {
	model := &Model{
		Table:         table,
		Columns:       columns,
		PrimaryKey:    primaryKey,
		Uuid:          false,
		AutoIncrement: true,
		Timestamps:    timestamps,
		SoftDelete:    softDelete,
	}

	model.handleTimestamp()
	model.handleSoftDelete()

	return model
}

type Model struct {
	Table         string
	Columns       []string
	PrimaryKey    string
	Uuid          bool
	AutoIncrement bool
	Timestamps    bool
	SoftDelete    bool
	CreatedAt     *time.Time `db:"createdAt"`
	UpdatedAt     *time.Time `db:"updatedAt"`
	DeletedAt     *time.Time `db:"deletedAt"`
}

// MUST OVERRIDE

func (m *Model) GeneratePK() {
	if m.IsUuid() {
		panic("Please override GeneratePK() method in your model")
	}
}

func (m *Model) Get(key string) interface{} {
	panic("Please override Get(key string) method in your model")
}

func (m *Model) Payload() map[string]interface{} {
	panic("Please override Payload() method in your model")
}

// LEAVE ALONE

func (m *Model) GetModel() *Model {
	return m
}

func (m *Model) GetTable() string {
	return m.Table
}

func (m *Model) GetColumns() []string {
	return m.Columns
}

func (m *Model) GetPK() string {
	return m.PrimaryKey
}

func (m *Model) GenerateId(v interface{}) {
	if !m.Uuid {
		return
	}

	pkField := strings.Title(m.GetPK())
	id := uuid.NewV4().String()
	reflect.ValueOf(v).Elem().FieldByName(pkField).Set(reflect.ValueOf(id))
}

func (m *Model) ColumnCount() int {
	return len(m.Columns)
}

func (m *Model) IsTimestamps() bool {
	return m.Timestamps
}

func (m *Model) IsSoftDelete() bool {
	return m.SoftDelete
}

func (m *Model) IsUuid() bool {
	return m.Uuid
}

func (m *Model) IsAutoIncrement() bool {
	return m.AutoIncrement
}

func (m *Model) MapPayload(v interface{}) map[string]interface{} {
	var results = make(map[string]interface{})
	classType := reflect.TypeOf(v)
	classValue := reflect.ValueOf(v)

	if classType.Kind() == reflect.Ptr {
		classType = classType.Elem()
	}
	if classValue.Kind() == reflect.Ptr {
		classValue = classValue.Elem()
	}

	for i := 0; i < classType.NumField(); i++ {
		if classType.Field(i).Name == "Model" {
			model := classValue.Field(i).Interface().(*Model)
			if model.IsTimestamps() {
				results["createdAt"] = model.CreatedAt
				results["updatedAt"] = model.UpdatedAt
			}
			if model.IsSoftDelete() {
				results["deletedAt"] = model.DeletedAt
			}
 		} else {
			results[m.toLowerFirst(classType.Field(i).Name)] = classValue.Field(i).Interface()
		}
	}
	return results
}

func (m *Model) Created() {
	now := time.Now()
	m.CreatedAt = &now
}

func (m *Model) Updated() {
	now := time.Now()
	m.UpdatedAt = &now
}

func (m *Model) Deleted() {
	now := time.Now()
	m.DeletedAt = &now
}

func (m *Model) Recovered() {
	m.DeletedAt = nil
}

func (m *Model) handleTimestamp() {
	if m.IsTimestamps() {
		m.Columns = append(m.Columns, "createdAt")
		m.Columns = append(m.Columns, "updatedAt")
	}
}

func (m *Model) handleSoftDelete() {
	if m.IsSoftDelete() {
		m.Columns = append(m.Columns, "deletedAt")
	}
}

func (m *Model) toLowerFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return str
}