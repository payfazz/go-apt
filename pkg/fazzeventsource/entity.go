package fazzeventsource

// EntityInterface is interface for entity
type EntityInterface interface {
	Apply(event Event) error
	GetId() string
	GetVersion() int
	GetUncommittedEvents() []Event
	ClearUncommittedEvents()
}

// Entity is struct for entity that is connected to the event
type Entity struct {
	Id                string `db:"id" json:"id"`
	Version           int    `db:"version" json:"version"`
	uncommittedEvents []Event
}

// Apply apply version and add event to entity
func (e *Entity) Apply(event Event) error {
	e.Version++
	e.uncommittedEvents = append(e.uncommittedEvents, event)
	return nil
}

// GetId return id of entity
func (e *Entity) GetId() string {
	return e.Id
}

// GetVersion return version of entity
func (e *Entity) GetVersion() int {
	return e.Version
}

// GetUncommittedEvents return uncommitted events
func (e *Entity) GetUncommittedEvents() []Event {
	return e.uncommittedEvents
}

// ClearUncommittedEvents clear uncommited events
func (e *Entity) ClearUncommittedEvents() {
	e.uncommittedEvents = nil
}
