package fazzdb

// Config is a struct that will be used to set default value for some parameter attribute
type Config struct {
	Limit  int
	Offset int
	Lock   Lock
}
