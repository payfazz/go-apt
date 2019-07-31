package eventstore

import (
	"context"
	"github.com/payfazz/go-apt/pkg/esfazz"
)

// EventStore is interface for event store
type EventStore interface {
	Save(ctx context.Context, event *esfazz.Event) error
	FindNotApplied(ctx context.Context, agg esfazz.Aggregate) ([]*esfazz.Event, error)
}
