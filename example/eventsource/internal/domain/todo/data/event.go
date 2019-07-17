package data

const (
	EVENT_TODO_CREATED = "todo.created"
	EVENT_TODO_UPDATED = "todo.updated"
	EVENT_TODO_DELETED = "todo.deleted"
)

// TodoCreated is data for todo created event
type TodoCreated struct {
	Id   string `json:"id"`
	Text string `json:"text"`
}

// TodoUpdated is data for todo updated event
type TodoUpdated struct {
	Id        string  `json:"id"`
	Text      *string `json:"text"`
	Completed *bool   `json:"completed"`
}

// TodoDeleted is data for todo deleted event
type TodoDeleted struct {
	Id string `json:"id"`
}
