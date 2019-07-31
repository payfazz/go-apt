package esrepo

import (
	"context"
	"github.com/payfazz/go-apt/pkg/esfazz"
)

// EventSourceRepository is interface for event repository
type EventSourceRepository interface {
	Save(ctx context.Context, event *esfazz.Event) error
	Find(ctx context.Context, id string) (interface{}, error)
}
