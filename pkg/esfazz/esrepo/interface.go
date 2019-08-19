package esrepo

import (
	"context"
	"github.com/payfazz/go-apt/pkg/esfazz"
)

// Repository is interface for event repository
type Repository interface {
	Save(ctx context.Context, events ...*esfazz.EventPayload) error
	Find(ctx context.Context, id string) (esfazz.Aggregate, error)
}

// EventListener is function that listen to event
type EventListener func(ctx context.Context, events ...*esfazz.Event) error
