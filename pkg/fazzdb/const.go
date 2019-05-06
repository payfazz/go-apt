package fazzdb

// Connector is a type that is used to connect the result of different condition
type Connector string

// Lock is a type that is used to set lock type used in query
type Lock string

// Operator is a type that is used to set operator relation in a condition
type Operator string

// OrderDirection is a type that is used to set order by direction for query
type OrderDirection string

// Aggregate is a type that is used to set aggregate used for query
type Aggregate string

// TrashStatus is a type that is used to set WITH_TRASH or NO_TRASH status
type TrashStatus bool

// DataType is a type that is used to set MigrationColumn type
type DataType string

// ReferenceAction is a type that is used to set MigrationReferences action
type ReferenceAction string

// TableCommand is a type that is used to set MigrationTable command
type MigrationCommand string

const (
	OP_EQUALS           Operator = "="
	OP_NOT_EQUALS       Operator = "!="
	OP_LESS_THAN        Operator = "<"
	OP_LESS_THAN_EQUALS Operator = "<="
	OP_MORE_THAN        Operator = ">"
	OP_MORE_THAN_EQUALS Operator = ">="
	OP_LIKE             Operator = "LIKE"
	OP_NOT_LIKE         Operator = "NOT LIKE"
	OP_IS_NULL          Operator = "IS NULL"
	OP_IS_NOT_NULL      Operator = "IS NOT NULL"
	OP_IN               Operator = "IN"
)

const (
	CO_AND  Connector = "AND"
	CO_OR   Connector = "OR"
	CO_NONE Connector = ""
)

const (
	LO_FOR_SHARE  Lock = "FOR SHARE"
	LO_FOR_UPDATE Lock = "FOR UPDATE"
	LO_NONE       Lock = ""
)

const (
	DIR_ASC  OrderDirection = "ASC"
	DIR_DESC OrderDirection = "DESC"
)

const (
	AG_NONE  Aggregate = ""
	AG_COUNT Aggregate = "COUNT"
	AG_SUM   Aggregate = "SUM"
	AG_AVG   Aggregate = "AVG"
	AG_MIN   Aggregate = "MIN"
	AG_MAX   Aggregate = "MAX"
)

const (
	WITH_TRASH TrashStatus = true
	NO_TRASH   TrashStatus = false
)

const (
	CREATED_AT = "created_at"
	UPDATED_AT = "updated_at"
	DELETED_AT = "deleted_at"
)

const (
	DT_UUID        DataType = "UUID"
	DT_STRING      DataType = "VARCHAR"
	DT_JSON        DataType = "JSON"
	DT_JSONB       DataType = "JSONB"
	DT_INT         DataType = "INTEGER"
	DT_SERIAL      DataType = "SERIAL"
	DT_BIGSERIAL   DataType = "BIGSERIAL"
	DT_BIGINT      DataType = "BIGINT"
	DT_BOOL        DataType = "BOOLEAN"
	DT_TEXT        DataType = "TEXT"
	DT_DOUBLE      DataType = "DOUBLE"
	DT_NUMERIC     DataType = "NUMERIC"
	DT_DECIMAL     DataType = "DECIMAL"
	DT_TIMESTAMP   DataType = "TIMESTAMP"
	DT_TIMESTAMPTZ DataType = "TIMESTAMPTZ"
	DT_NONE        DataType = "NONE"
)

const (
	RA_NO_ACTION ReferenceAction = "NO ACTION"
	RA_CASCADE   ReferenceAction = "CASCADE"
	RA_RESTRICT  ReferenceAction = "RESTRICT"
)

const (
	MC_CREATE MigrationCommand = "CREATE"
	MC_ADD    MigrationCommand = "ADD"
	MC_ALTER  MigrationCommand = "ALTER"
	MC_DROP   MigrationCommand = "DROP"
	MC_RENAME MigrationCommand = "RENAME"
)

const (
	META_APP_ID  = "APP_ID"
	META_VERSION = "VERSION"
	META_TABLE   = "fazz_metas"
)

var DEFAULT_QUERY_CONFIG = Config{
	Limit:  0,
	Offset: 0,
	Lock:   LO_NONE,
}
