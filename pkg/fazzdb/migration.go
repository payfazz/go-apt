package fazzdb

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/payfazz/go-apt/pkg/fazzcommon/formatter"
	"log"
)

// FazzMetaModel is a constructor of FazzMeta model
func FazzMetaModel() *FazzMeta {
	return &FazzMeta{
		Model: *PlainModel(
			"fazz_metas",
			[]Column{
				Col("key"),
				Col("value"),
			},
			"key",
			false,
			false,
		),
	}
}

// FazzMeta is a struct that will contain key and value from fazzdb meta table
type FazzMeta struct {
	Model
	Key   string `db:"key"`
	Value string `db:"value"`
}

// Get is overriden function from Model
func (fm *FazzMeta) Get(key string) interface{} {
	return fm.Payload()[key]
}

// Payload is overriden function from Model
func (fm *FazzMeta) Payload() map[string]interface{} {
	return fm.MapPayload(fm)
}

// Migrate is a constructor of Migration struct that will run creation and seed of meta table and other migration
func Migrate(db *sqlx.DB, appId string, versions ...MigrationVersion) {
	m := &Migration{
		Versions: versions,
	}

	tx, _ := db.Beginx()
	query := QueryTx(tx, DEFAULT_QUERY_CONFIG)

	m.RunMeta(query, appId)
	if m.isRightApp(query, appId) {
		m.Run(query)
	} else {
		panic("migrating to wrong app")
	}

	_ = query.Tx.Commit()
}

// Migration is a struct that contains versions of migration used in the app
type Migration struct {
	Versions []MigrationVersion
}

// RunMeta is a function that will create and seed fazzdb meta table
func (m *Migration) RunMeta(query *Query, appId string) {
	metaMigration := MigrationVersion{
		Tables: []*MigrationTable{
			CreateTable(META_TABLE, func(table *MigrationTable) {
				table.Field(CreateText("key").Primary())
				table.Field(CreateText("value"))
			}),
		},
	}

	log.Printf("Creating meta table")
	metaMigration.Run(query, false)

	if !m.isMetaExist(query) {
		m.seedMetaAppId(query, appId)
		m.seedMetaVersion(query)
	}
}

// Run is a function that will run all migration tables and enums from Versions in Migration
func (m *Migration) Run(query *Query) {
	metaVersion := m.metaVersion(query)
	appVersion := m.appVersion()

	log.Printf("Meta version: %d", metaVersion)
	log.Printf("App version: %d", appVersion)

	if appVersion < metaVersion {
		panic("meta version is bigger than app version")
	} else if appVersion > metaVersion {
		for index, v := range m.Versions {
			if index >= metaVersion {
				log.Printf("Running migration version %d", metaVersion + 1)
				v.Run(query, true)

				metaVersion++
				m.incrementMetaVersion(query, metaVersion)
				log.Printf("Migration version %d done!", metaVersion)
			}
		}
		_ = query.Tx.Commit()
	} else {
		log.Println("Same meta and app version, doing nothing!")
	}
}

// incrementMetaVersion is a function that will add meta version by 1
func (m *Migration) incrementMetaVersion(query *Query, metaVersion int) {
	versionModel := FazzMetaModel()
	versionModel.Key = META_VERSION
	versionModel.Value = formatter.IntegerToString(metaVersion)

	_, err := query.Use(versionModel).Update()
	if nil != err {
		_ = query.Tx.Rollback()
		panic(err)
	}
}

// metaVersion is a function that will get and lock the version row from db until commit
func (m *Migration) metaVersion(query *Query) int {
	row, err := query.Use(FazzMetaModel()).
		Where("key", META_VERSION).
		WithLock(LO_FOR_SHARE).
		First()

	if nil != err {
		_ = query.Tx.Rollback()
		panic(err)
	}

	version := row.(FazzMeta)
	return formatter.StringToInteger(version.Value)
}

// appVersion is a function that will get current migration version from the code
func (m *Migration) appVersion() int {
	return len(m.Versions)
}

// isMetaExist is a function that will check if APP_ID and VERSION row is exist in database
func (m *Migration) isMetaExist(query *Query) bool {
	count, err := query.Use(FazzMetaModel()).
		Where("key", META_APP_ID).
		OrWhere("key", META_VERSION).
		Count()

	if nil != err {
		_ = query.Tx.Rollback()
		panic(err)
	}

	if int64(*count) == 2 {
		return true
	}

	return false
}

// isRightApp is a function that will check if current database have the right APP_ID before migrating
func (m *Migration) isRightApp(query *Query, appId string) bool {
	_, err := query.Use(FazzMetaModel()).
		Where("key", META_APP_ID).
		Where("value", appId).
		First()

	return nil == err
}

// seedMetaAppId is a function that will seed APP_ID to new meta table
func (m *Migration) seedMetaAppId(query *Query, appId string) {
	appIdModel := FazzMetaModel()
	appIdModel.Key = META_APP_ID
	appIdModel.Value = appId

	log.Printf("Seeding meta app id")
	_, err := query.Use(appIdModel).InsertOnConflict(true)
	if nil != err {
		_ = query.Tx.Rollback()
		panic(err)
	}
}

// seedMetaVersion is a function that will seed VERSION to new meta table
func (m *Migration) seedMetaVersion(query *Query) {
	versionModel := FazzMetaModel()
	versionModel.Key = META_VERSION
	versionModel.Value = "0"

	log.Printf("Seeding meta version")
	_, err := query.Use(versionModel).InsertOnConflict(true)
	if nil != err {
		_ = query.Tx.Rollback()
		panic(err)
	}
}

// CreateTable is a wrapper to table constructor for creating migration table
func CreateTable(name string, detail func(table *MigrationTable)) *MigrationTable {
	return newTable(name, MC_CREATE, detail)
}

// AlterTable is a wrapper to table constructor for altering migration table
func AlterTable(name string, detail func(table *MigrationTable)) *MigrationTable {
	return newTable(name, MC_ALTER, detail)
}

// DropTable is a wrapper to table constructor for dropping migration table
func DropTable(name string, detail func(table *MigrationTable)) *MigrationTable {
	return newTable(name, MC_DROP, detail)
}

// newTable is a constructor for MigrationTable that will check and parse primary key
func newTable(name string, command MigrationCommand, detail func(table *MigrationTable)) *MigrationTable {
	table := &MigrationTable{
		command: command,
		name:    name,
	}
	detail(table)
	table.parsePrimaryKey()
	table.validate()
	return table
}

// MigrationTable is a struct that is used to store information about migration table
type MigrationTable struct {
	command     MigrationCommand
	name        string
	columns     []*MigrationColumn
	references  []*MigrationReference
	primaryKeys []string
}

// Field is a function that will append new MigrationColumn to columns
func (mt *MigrationTable) Field(column *MigrationColumn) {
	mt.columns = append(mt.columns, column)
}

// Reference is a function that will append new MigrationReference to references
func (mt *MigrationTable) Reference(reference *MigrationReference) {
	mt.references = append(mt.references, reference)
}

// Timestamps is a function that will add createdAt and updatedAt column automatically
func (mt *MigrationTable) Timestamps() {
	mt.Field(CreateTimestamp("createdAt").Nullable())
	mt.Field(CreateTimestamp("updatedAt").Nullable())
}

// TimestampsTz is a function that will add createdAt and updatedAt column automatically with timezone
func (mt *MigrationTable) TimestampsTz(timezone int) {
	mt.Field(CreateTimestampTz("createdAt", timezone).Nullable())
	mt.Field(CreateTimestampTz("updatedAt", timezone).Nullable())
}

// SoftDelete is a function that will add deletedAt column automatically
func (mt *MigrationTable) SoftDelete() {
	mt.Field(CreateTimestamp("deletedAt").Nullable())
}

// SoftDeleteTz is a function that will add deletedAt column automatically with timezone
func (mt *MigrationTable) SoftDeleteTz(timezone int) {
	mt.Field(CreateTimestampTz("deletedAt", timezone).Nullable())
}

// parsePrimaryKey is a function that will append all primary key column into primaryKeys
func (mt *MigrationTable) parsePrimaryKey() {
	for _, column := range mt.columns {
		if column.primaryKey {
			mt.primaryKeys = append(mt.primaryKeys, column.name)
		}
	}
}

// validate is a function that will validate if some columns action is allowed
func (mt *MigrationTable) validate() {
	if mt.command == MC_RENAME || mt.command == MC_ADD {
		panic("invalid migration command")
	}

	if len(mt.columns) == 0 && mt.command != MC_DROP {
		panic("no columns on alter table migration")
	}

	for _, column := range mt.columns {
		if column.command == MC_RENAME && len(mt.columns) > 1 {
			panic("can't combine rename with other alter table query, please separate it")
		}
	}
}

// CreateUuid is a function to create MigrationColumn of DataType Uuid
func CreateUuid(name string) *MigrationColumn {
	return createColumn(name, DT_UUID)
}

// CreateInteger is a function to create MigrationColumn of DataType Integer
func CreateInteger(name string) *MigrationColumn {
	return createColumn(name, DT_INT)
}

// CreateBigInteger is a function to create MigrationColumn of DataType BigInteger
func CreateBigInteger(name string) *MigrationColumn {
	return createColumn(name, DT_BIGINT)
}

// CreateSerial is a function to create MigrationColumn of DataType Serial
func CreateSerial(name string) *MigrationColumn {
	return createColumn(name, DT_SERIAL)
}

// CreateString is a function to create MigrationColumn of DataType Varchar
func CreateString(name string) *MigrationColumn {
	return createColumn(name, DT_STRING)
}

// CreateBoolean is a function to create MigrationColumn of DataType Boolean
func CreateBoolean(name string) *MigrationColumn {
	return createColumn(name, DT_BOOL)
}

// CreateText is a function to create MigrationColumn of DataType Text
func CreateText(name string) *MigrationColumn {
	return createColumn(name, DT_TEXT)
}

// CreateDouble is a function to create MigrationColumn of DataType Double
func CreateDouble(name string) *MigrationColumn {
	return createColumn(name, DT_DOUBLE)
}

// CreateNumeric is a function to create MigrationColumn of DataType Numeric
func CreateNumeric(name string, precision int, scale int) *MigrationColumn {
	return createColumn(name, DT_NUMERIC, precision, scale)
}

// CreateDecimal is a function to create MigrationColumn of DataType Decimal
func CreateDecimal(name string, precision int, scale int) *MigrationColumn {
	return createColumn(name, DT_DECIMAL, precision, scale)
}

// CreateTimestamp is a function to create MigrationColumn of DataType Timestamp
func CreateTimestamp(name string) *MigrationColumn {
	return createColumn(name, DT_TIMESTAMP)
}

// CreateTimestampTz is a function to create MigrationColumn of DataType TimestampTz
func CreateTimestampTz(name string, timezone int) *MigrationColumn {
	return createColumn(name, DT_TIMESTAMPTZ, timezone)
}

// CreateEnum is a function to create MigrationColumn of DataType Enum
func CreateEnum(name string, enum MigrationEnum) *MigrationColumn {
	return createColumn(name, enum.GetDataType())
}

// AddUuid is a function to add MigrationColumn of DataType Uuid
func AddUuid(name string) *MigrationColumn {
	return addColumn(name, DT_UUID)
}

// AddInteger is a function to add MigrationColumn of DataType Integer
func AddInteger(name string) *MigrationColumn {
	return addColumn(name, DT_INT)
}

// AddBigInteger is a function to add MigrationColumn of DataType BigInteger
func AddBigInteger(name string) *MigrationColumn {
	return addColumn(name, DT_BIGINT)
}

// AddString is a function to add MigrationColumn of DataType String
func AddString(name string) *MigrationColumn {
	return addColumn(name, DT_STRING)
}

// AddBoolean is a function to add MigrationColumn of DataType Boolean
func AddBoolean(name string) *MigrationColumn {
	return addColumn(name, DT_BOOL)
}

// AddText is a function to add MigrationColumn of DataType Text
func AddText(name string) *MigrationColumn {
	return addColumn(name, DT_TEXT)
}

// AddDouble is a function to add MigrationColumn of DataType Double
func AddDouble(name string) *MigrationColumn {
	return addColumn(name, DT_DOUBLE)
}

// AddNumeric is a function to add MigrationColumn of DataType Numeric
func AddNumeric(name string, precision int, scale int) *MigrationColumn {
	return addColumn(name, DT_NUMERIC, precision, scale)
}

// AddDecimal is a function to add MigrationColumn of DataType Decimal
func AddDecimal(name string, precision int, scale int) *MigrationColumn {
	return addColumn(name, DT_DECIMAL, precision, scale)
}

// AddTimestamp is a function to add MigrationColumn of DataType Timestamp
func AddTimestamp(name string) *MigrationColumn {
	return addColumn(name, DT_TIMESTAMP).Nullable()
}

// AddTimestampTz is a function to add MigrationColumn of DataType TimestampTz
func AddTimestampTz(name string, timezone int) *MigrationColumn {
	return addColumn(name, DT_TIMESTAMPTZ, timezone)
}

// AddEnum is a function to add MigrationColumn of DataType Enum
func AddEnum(name string, enum MigrationEnum) *MigrationColumn {
	return addColumn(name, enum.GetDataType())
}

// AlterUuid is a function to alter MigrationColumn to DataType Uuid
func AlterUuid(name string) *MigrationColumn {
	return alterColumn(name, DT_UUID)
}

// AlterInteger is a function to alter MigrationColumn to DataType Integer
func AlterInteger(name string) *MigrationColumn {
	return alterColumn(name, DT_INT)
}

// AlterBigInteger is a function to alter MigrationColumn to DataType BigInteger
func AlterBigInteger(name string) *MigrationColumn {
	return alterColumn(name, DT_BIGINT)
}

// AlterString is a function to alter MigrationColumn to DataType String
func AlterString(name string) *MigrationColumn {
	return alterColumn(name, DT_STRING)
}

// AlterBoolean is a function to alter MigrationColumn to DataType Boolean
func AlterBoolean(name string) *MigrationColumn {
	return alterColumn(name, DT_BOOL)
}

// AlterText is a function to alter MigrationColumn to DataType Text
func AlterText(name string) *MigrationColumn {
	return alterColumn(name, DT_TEXT)
}

// AlterDouble is a function to alter MigrationColumn to DataType Double
func AlterDouble(name string) *MigrationColumn {
	return alterColumn(name, DT_DOUBLE)
}

// AlterNumeric is a function to alter MigrationColumn to DataType Numeric
func AlterNumeric(name string, precision int, scale int) *MigrationColumn {
	return alterColumn(name, DT_NUMERIC, precision, scale)
}

// AlterDecimal is a function to alter MigrationColumn to DataType Decimal
func AlterDecimal(name string, precision int, scale int) *MigrationColumn {
	return alterColumn(name, DT_DECIMAL, precision, scale)
}

// AlterTimestamp is a function to alter MigrationColumn to DataType Timestamp
func AlterTimestamp(name string) *MigrationColumn {
	return alterColumn(name, DT_TIMESTAMP)
}

// AlterTimestampTz is a function to alter MigrationColumn to DataType TimestampTz
func AlterTimestampTz(name string, timezone int) *MigrationColumn {
	return alterColumn(name, DT_TIMESTAMPTZ, timezone)
}

// AlterEnum is a function to alter MigrationColumn to DataType Enum
func AlterEnum(name string, enum MigrationEnum) *MigrationColumn {
	return alterColumn(name, enum.GetDataType())
}

// DropColumn is a function to drop MigrationColumn
func DropColumn(name string) *MigrationColumn {
	return newColumn("", name, DT_NONE, MC_DROP)
}

// RenameColumn is a function to rename MigrationColumn
func RenameColumn(previousName string, name string) *MigrationColumn {
	return newColumn(previousName, name, DT_NONE, MC_RENAME)
}

// createColumn is a wrapper to the MigrationColumn constructor for creating column
func createColumn(name string, dataType DataType, args ...int) *MigrationColumn {
	return newColumn("", name, dataType, MC_CREATE, args...)
}

// addColumn is a wrapper to the MigrationColumn constructor for adding column
func addColumn(name string, dataType DataType, args ...int) *MigrationColumn {
	return newColumn("", name, dataType, MC_ADD, args...)
}

// alterColumn is a wrapper to the MigrationColumn constructor for altering column
func alterColumn(name string, dataType DataType, args ...int) *MigrationColumn {
	return newColumn("", name, dataType, MC_ALTER, args...)
}

// newColumn is a constructor of MigrationColumn with default values
func newColumn(previousName string, name string, dataType DataType, command MigrationCommand, args ...int) *MigrationColumn {
	return &MigrationColumn{
		command:      command,
		previousName: previousName,
		name:         name,
		dataType:     dataType,
		typeArgs:     args,
		primaryKey:   false,
		unique:       false,
		nullable:     false,
	}
}

// MigrationColumn is a struct that is used to store information about migration column
type MigrationColumn struct {
	command      MigrationCommand
	previousName string
	name         string
	dataType     DataType
	typeArgs     []int
	primaryKey   bool
	unique       bool
	nullable     bool
	defaultValue string
}

// Primary is a function to set primaryKey flag to true in MigrationColumn
func (mc *MigrationColumn) Primary() *MigrationColumn {
	mc.primaryKey = true
	return mc
}

// Unique is a function to set unique flag to true in MigrationColumn
func (mc *MigrationColumn) Unique() *MigrationColumn {
	mc.unique = true
	return mc
}

// Nullable is a function to set nullable flag to true in MigrationColumn
func (mc *MigrationColumn) Nullable() *MigrationColumn {
	if mc.primaryKey {
		panic("primary key cannot be nullable")
	}
	mc.nullable = true
	return mc
}

// Default is a function to set Default value in MigrationColumn
func (mc *MigrationColumn) Default(value string) *MigrationColumn {
	mc.defaultValue = value
	return mc
}

// MigrationReference is a struct that is used to store information about migration references
type MigrationReference struct {
	key        string
	otherKey   string
	otherTable string
	onUpdate   ReferenceAction
	onDelete   ReferenceAction
}

// Foreign is a constructor to create MigrationReference for foreign key constraint
func Foreign(key string) *MigrationReference {
	reference := &MigrationReference{
		key:      key,
		onUpdate: RA_NO_ACTION,
		onDelete: RA_NO_ACTION,
	}
	return reference
}

// Reference is a function to set key reference in other table
func (mr *MigrationReference) Reference(key string) *MigrationReference {
	mr.otherKey = key
	return mr
}

// On is a function to set which table to reference
func (mr *MigrationReference) On(table string) *MigrationReference {
	mr.otherTable = table
	return mr
}

// OnUpdate is a function to set reference action when data is updated
func (mr *MigrationReference) OnUpdate(action ReferenceAction) *MigrationReference {
	mr.onUpdate = action
	return mr
}

// OnDelete is a function to set reference action when data is deleted
func (mr *MigrationReference) OnDelete(action ReferenceAction) *MigrationReference {
	mr.onDelete = action
	return mr
}

// MigrationEnum is a struct that is used to store information about migration enum
type MigrationEnum struct {
	name   string
	values []string
}

// NewEnum is a constructor to MigrationEnum
func NewEnum(name string, values ...string) *MigrationEnum {
	return &MigrationEnum{
		name:   fmt.Sprintf("%s_enum", name),
		values: values,
	}
}

// GetType is a function to return enum name with public prefix
func (me *MigrationEnum) GetType() string {
	return fmt.Sprintf(`"public"."%s"`, me.name)
}

// GetDataType is a function to return enum data type for column use
func (me *MigrationEnum) GetDataType() DataType {
	return DataType(me.GetType())
}

// MigrationVersion is a struct that is used to store information about one version of migration
type MigrationVersion struct {
	Tables []*MigrationTable
	Enums  []*MigrationEnum
}

// Run is a function that will run all tables and enums command in a MigrationVersion
func (mv *MigrationVersion) Run(query *Query, autoDrop bool) {
	builder := NewBuilder()

	queryString := ``
	if autoDrop {
		for i := len(mv.Tables) - 1; i >= 0; i-- {
			if mv.Tables[i].command == MC_CREATE {
				queryString = fmt.Sprintf(`%s %s`, queryString, builder.BuildDropTable(mv.Tables[i]))
			}
		}
		for i := len(mv.Enums) - 1; i >= 0; i-- {
			queryString = fmt.Sprintf(`%s %s`, queryString, builder.BuildDropEnum(mv.Enums[i]))
		}
	}

	for _, enum := range mv.Enums {
		queryString = fmt.Sprintf(`%s %s`, queryString, builder.BuildCreateEnum(enum))
	}

	for _, table := range mv.Tables {
		if table.command == MC_CREATE {
			queryString = fmt.Sprintf(`%s %s`, queryString, builder.BuildCreateTable(table))
		} else if table.command == MC_ALTER {
			queryString = fmt.Sprintf(`%s %s`, queryString, builder.BuildAlterTable(table))
		} else if table.command == MC_DROP {
			queryString = fmt.Sprintf(`%s %s`, queryString, builder.BuildDropTable(table))
		}
	}

	_, err := query.RawExec(queryString)
	if nil != err {
		_ = query.Tx.Rollback()
		panic(err)
	}
}