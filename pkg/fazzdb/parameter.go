package fazzdb

import (
	"fmt"
	"log"
	"time"
)

// NewParameter is a constructor that will return Parameter instance
func NewParameter(config Config) *Parameter {
	return &Parameter{
		Values:          make(map[string]interface{}, 0),
		Offset:          config.Offset,
		Limit:           config.Limit,
		Lock:            config.Lock,
		DevelopmentMode: config.DevelopmentMode,
	}
}

// Parameter is a struct that is used to contain all configuration of your query
type Parameter struct {
	Conditions      []Condition
	Havings         []Condition
	Values          map[string]interface{}
	Columns         []Column
	Orders          []Order
	Groups          []Column
	Lock            Lock
	Limit           int
	Offset          int
	DevelopmentMode bool
}

// appendGroupConditions is a function that is used to append multiple conditions as one new Conditions
func (p *Parameter) appendGroupConditions(param *Parameter, connector Connector) *Parameter {
	return p.appendConditionFromParameter(param, connector)
}

// appendGroupHavings is a function that is used to append multiple havings as one new Havings
func (p *Parameter) appendGroupHavings(param *Parameter, connector Connector) *Parameter {
	return p.appendConditionFromParameter(param, connector)
}

// appendConditionFromParameter is a function that append both Conditions and Havings if available
func (p *Parameter) appendConditionFromParameter(param *Parameter, connector Connector) *Parameter {
	// Append Condition
	if len(param.Conditions) > 0 {
		conditionParent := Condition{
			Connector: connector,
		}
		for i, condition := range param.Conditions {
			if i == 0 {
				condition.Connector = CO_NONE
			}
			conditionParent.Conditions = append(conditionParent.Conditions, condition)
		}
		p.Conditions = append(p.Conditions, conditionParent)
	}

	// Append Having
	if len(param.Havings) > 0 {
		havingParent := Condition{
			Connector: connector,
		}
		for i, having := range param.Havings {
			if i == 0 {
				having.Connector = CO_NONE
			}
			havingParent.Conditions = append(havingParent.Conditions, having)
		}
		p.Havings = append(p.Havings, havingParent)
	}

	// Append Values
	for i, value := range param.Values {
		p.Values[i] = value
	}

	return p
}

// appendCondition is a function to append single condition to Conditions attribute
func (p *Parameter) appendCondition(
	connector Connector,
	field Column,
	operator Operator,
	value interface{},
) *Parameter {
	prefix := p.getPrefix(field.Key)
	p.Conditions = append(p.Conditions, Condition{
		Field:     field,
		Operator:  operator,
		Connector: connector,
		Prefix:    prefix,
	})

	if operator == OP_IS_NOT_NULL || operator == OP_IS_NULL {
		return p
	}

	if _, ok := p.Values[prefix]; ok {
		log.Printf("[ALERT] Collision on condition prefix: %s", prefix)
	}

	p.Values[prefix] = value
	return p
}

// appendHaving is a function to append single having to Havings attribute
func (p *Parameter) appendHaving(
	connector Connector,
	field Column,
	operator Operator,
	value interface{},
) *Parameter {
	prefix := p.getPrefix(field.Key)
	p.Havings = append(p.Havings, Condition{
		Field:     field,
		Operator:  operator,
		Connector: connector,
		Prefix:    prefix,
	})

	if _, ok := p.Values[prefix]; ok {
		log.Printf("[ALERT] Collision on having prefix: %s", prefix)
	}

	p.Values[prefix] = value
	return p
}

// appendOrderBy is a function to append order by column to Orders attribute
func (p *Parameter) appendOrderBy(table string, field Column, direction OrderDirection, nullsLast bool) *Parameter {
	p.Orders = append(p.Orders, Order{
		Field:     field,
		Direction: direction,
		NullsLast: nullsLast,
	})
	return p
}

// appendGroupBy is a function to append group by column to Groups attribute
func (p *Parameter) appendGroupBy(column Column) *Parameter {
	p.Groups = append(p.Groups, column)
	return p
}

// setColumns is a function to assign Columns attribute with given parameter
func (p *Parameter) setColumns(columns []Column) *Parameter {
	p.Columns = columns
	return p
}

// setLock is a function to add lock type to Lock attribute
func (p *Parameter) setLock(lock Lock) *Parameter {
	p.Lock = lock
	return p
}

// setLimit is a function to set Limit attribute
func (p *Parameter) setLimit(limit int) *Parameter {
	p.Limit = limit
	return p
}

// setOffset is a function to set Offset attribute
func (p *Parameter) setOffset(offset int) *Parameter {
	p.Offset = offset
	return p
}

// setDevelopmentMode is a function to set DevelopmentMode attribute
func (p *Parameter) setDevelopmentMode(developmentMode bool) *Parameter {
	p.DevelopmentMode = developmentMode
	return p
}

// getPrefix is a function to generate prefix for condition arguments using nanoseconds
func (p *Parameter) getPrefix(key string) string {
	postfix := time.Now().UnixNano() / int64(time.Microsecond) % 100000
	return fmt.Sprintf("%d.%s.%d", len(p.Conditions), key, postfix)
}
