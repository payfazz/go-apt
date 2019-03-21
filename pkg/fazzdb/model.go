package fazzdb

import (
	"github.com/satori/go.uuid"
	"reflect"
	"strings"
	"unicode"
)

type ModelInterface interface {
	SetModel(v *Model)
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
}

type Model struct {
	Table         string
	Columns       []string
	PrimaryKey    string
	Uuid          bool
	AutoIncrement bool
	Timestamps    bool
	SoftDelete    bool
}

// MUST OVERRIDE

func (m *Model) GeneratePK() {
	panic("Please override GeneratePK() method in your model")
}

func (m *Model) Get(key string) interface{} {
	panic("Please override Get(key string) method in your model")
	return nil
}

func (m *Model) Payload() map[string]interface{} {
	panic("Please override Payload() method in your model")
	return make(map[string]interface{})
}

// LEAVE ALONE

func (m *Model) SetModel(v *Model) {
	m = v
}

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
		if classType.Field(i).Name != "Model" {
			results[m.toLowerFirst(classType.Field(i).Name)] = classValue.Field(i).Interface()
		}
	}
	return results
}

func (m *Model) toLowerFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return str
}

func UuidModel(table string, columns []string, primaryKey string, timestamps bool, softDelete bool) *Model {
	return &Model{
		Table:         table,
		Columns:       columns,
		PrimaryKey:    primaryKey,
		Uuid:          true,
		AutoIncrement: false,
		Timestamps:    timestamps,
		SoftDelete:    softDelete,
	}
}

func PlainModel(table string, columns []string, primaryKey string, timestamps bool, softDelete bool) *Model {
	return &Model{
		Table:         table,
		Columns:       columns,
		PrimaryKey:    primaryKey,
		Uuid:          false,
		AutoIncrement: false,
		Timestamps:    timestamps,
		SoftDelete:    softDelete,
	}
}

func AutoIncrementModel(table string, columns []string, primaryKey string, timestamps bool, softDelete bool) *Model {
	return &Model{
		Table:         table,
		Columns:       columns,
		PrimaryKey:    primaryKey,
		Uuid:          false,
		AutoIncrement: true,
		Timestamps:    timestamps,
		SoftDelete:    softDelete,
	}
}
