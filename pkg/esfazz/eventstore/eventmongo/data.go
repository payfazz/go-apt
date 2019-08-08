package eventmongo

import "github.com/payfazz/go-apt/pkg/esfazz"

type eventLog struct {
	Type      string
	Aggregate esfazz.BaseAggregate
	Data      map[string]interface{}
}
