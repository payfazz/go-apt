package fazzdb

import (
	"reflect"
	"unicode"
)

type ModelInterface interface {
	SetModel(v *Model)
	GetModel() *Model
	GetTable() string
	GetColumns() []string
	GetPrimaryKey() string
	Get(key string) interface{}
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

func (m *Model) GetPrimaryKey() string {
	return m.PrimaryKey
}

func (m *Model) Get(key string) interface{} {
	return nil
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

func (m *Model) Payload() map[string]interface{} {
	return make(map[string]interface{})
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

func NewUuidModel(table string, columns []string, primaryKey string, timestamps bool, softDelete bool) *Model {
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

func NewPlainModel(table string, columns []string, primaryKey string, timestamps bool, softDelete bool) *Model {
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

func NewModel(table string, columns []string, primaryKey string, timestamps bool, softDelete bool) *Model {
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