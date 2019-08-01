package esfazz

// Aggregate is interface for aggregate that can apply event
type Aggregate interface {
	GetId() string
	GetVersion() int
	Apply(events ...*Event) error
}

// AggregateFactory is function type that create aggregate
type AggregateFactory func(id string) Aggregate

// BaseAggregate is base aggregate
type BaseAggregate struct {
	Id      string `json:"id"`
	Version int    `json:"version"`
}

// GetId return id of base aggregate
func (a *BaseAggregate) GetId() string {
	return a.Id
}

// GetVersion return version of base aggregate
func (a *BaseAggregate) GetVersion() int {
	return a.Version
}

// Apply implemented in base aggregate to for aggregate interface
func (a *BaseAggregate) Apply(events ...*Event) error {
	return nil
}
