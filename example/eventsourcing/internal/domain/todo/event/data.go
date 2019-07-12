package event

// TodoCreatedData is data for todo created event
type TodoCreatedData struct {
	Id   string `json:"id"`
	Text string `json:"text"`
}

// TodoUpdatedData is data for todo updated event
type TodoUpdatedData struct {
	Id        string `json:"id"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}

// TodoDeletedData is data for todo deleted event
type TodoDeletedData struct {
	Id string `json:"id"`
}
