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
	CO_AND   Connector = "AND"
	CO_OR    Connector = "OR"
	CO_EMPTY Connector = ""
)

const (
	LO_FOR_SHARE  Lock = "FOR SHARE"
	LO_FOR_UPDATE Lock = "FOR UPDATE"
	LO_EMPTY      Lock = ""
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
	CREATED_AT = "createdAt"
	UPDATED_AT = "updatedAt"
	DELETED_AT = "deletedAt"
)
