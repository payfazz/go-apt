package esfazz

type Aggregate interface {
	GetId() string
	GetVersion() int
}

type BaseAggregate struct {
	Id      string `json:"id"`
	Version int    `json:"version"`
}

func (a *BaseAggregate) GetId() string {
	return a.Id
}

func (a *BaseAggregate) GetVersion() int {
	return a.Version
}
