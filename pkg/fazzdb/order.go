package fazzdb

// Order is a struct that is used to contain order by attributes
type Order struct {
	Table     string
	Key       string
	Direction OrderDirection
}
