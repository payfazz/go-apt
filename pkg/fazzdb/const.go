package fazzdb

type Connector string
type Lock string
type Operator string
type OrderDirection string
type Aggregate string
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
