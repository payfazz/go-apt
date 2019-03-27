package config

import "github.com/payfazz/go-apt/pkg/fazzdb"

var Db = fazzdb.Config{
	Limit:  0,
	Offset: 0,
	Lock:   fazzdb.LO_EMPTY,
}
