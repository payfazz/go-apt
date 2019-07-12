package query

// Todo is a data that is returned for todo service
type Todo struct {
	Id        string `json:"id" db:"id"`
	Text      string `json:"text" db:"text"`
	Completed bool   `json:"completed" db:"completed"`
}
