package fazzdb

import (
	"db/fazzdb/fazzorder"
	"db/fazzdb/fazzspec"
)

type Parameter struct {
	Conditions []fazzspec.Condition
	Values     map[string]interface{}
	Orders     []fazzorder.Order
	Lock       fazzspec.Lock
	Limit      int
	Offset     int
}

func (p *Parameter) appendCondition(
	table string,
	connector fazzspec.Connector,
	key string,
	operator fazzspec.Operator,
	value interface{},
) *Parameter {
	p.Conditions = append(p.Conditions, fazzspec.Condition{
		Table:     table,
		Key:       key,
		Operator:  operator,
		Connector: connector,
	})
	p.Values[key] = value
	return p
}

func (p *Parameter) appendOrder(table string, key string, direction fazzorder.OrderDirection) *Parameter {
	p.Orders = append(p.Orders, fazzorder.Order{
		Table:     table,
		Key:       key,
		Direction: direction,
	})
	return p
}

func (p *Parameter) setLock(lock fazzspec.Lock) *Parameter {
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

func NewParameter() *Parameter {
	return &Parameter{
		Values: make(map[string]interface{}, 0),
	}
}
