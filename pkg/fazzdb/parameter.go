package fazzdb

import (
	"strconv"
	"time"
)

type Parameter struct {
	Conditions []Condition
	Values     map[string]interface{}
	Orders     []Order
	Groups     []string
	Lock       Lock
	Limit      int
	Offset     int
}

func (p *Parameter) appendGroupConditions(param *Parameter, connector Connector) *Parameter {
	// Append Condition
	parent := Condition{
		Connector: connector,
	}
	for i, condition := range param.Conditions {
		if i == 0 {
			condition.Connector = CO_EMPTY
		}
		parent.Conditions = append(parent.Conditions, condition)
	}
	p.Conditions = append(p.Conditions, parent)

	// Append Values
	for i, value := range param.Values {
		p.Values[i] = value
	}

	return p
}

func (p *Parameter) appendCondition(
	table string,
	connector Connector,
	key string,
	operator Operator,
	value interface{},
) *Parameter {
	prefix := p.getPrefix()
	p.Conditions = append(p.Conditions, Condition{
		Table:     table,
		Key:       key,
		Operator:  operator,
		Connector: connector,
		Prefix:    prefix,
	})
	p.Values[prefix] = value
	return p
}

func (p *Parameter) appendOrderBy(table string, key string, direction OrderDirection) *Parameter {
	p.Orders = append(p.Orders, Order{
		Table:     table,
		Key:       key,
		Direction: direction,
	})
	return p
}

func (p *Parameter) appendGroupBy(column string) *Parameter {
	p.Groups = append(p.Groups, column)
	return p
}

func (p *Parameter) setLock(lock Lock) *Parameter {
	p.Lock = lock
	return p
}

func (p *Parameter) setLimit(limit int) *Parameter {
	p.Limit = limit
	return p
}

func (p *Parameter) setOffset(offset int) *Parameter {
	p.Offset = offset
	return p
}

func (p *Parameter) getPrefix() string {
	prefix := (time.Now().UnixNano() / int64(time.Microsecond)) % int64(100000000)
	return strconv.Itoa(int(prefix))
}

func NewParameter(config Config) *Parameter {
	return &Parameter{
		Values: make(map[string]interface{}, 0),
		Offset: config.Offset,
		Limit:  config.Limit,
		Lock:   config.Lock,
	}
}
