package fazzspec

type Operator string
type Connector string
type Lock string

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
	OP_IS_NOT_NULL		Operator = "IS NOT NULL"
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
