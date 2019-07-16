package data

// PayloadCreateTodo is a struct for payload when creating todo
type PayloadCreateTodo struct {
	Text string `json:"text"`
}

// PayloadUpdateTodo is a struct for payload when updating todo
type PayloadUpdateTodo struct {
	Id        string `json:"id"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}
