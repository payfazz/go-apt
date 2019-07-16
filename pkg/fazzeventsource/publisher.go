package fazzeventsource

import (
	"context"
	"encoding/json"
	"github.com/payfazz/go-apt/pkg/fazzpubsub"
)

type EventPublisher interface {
	Publish(ctx context.Context, event *Event) error
}

type eventPublisher struct {
	pubsub    fazzpubsub.PubSub
	topicName string
}

// Publish do publishing to event pubsub
func (e *eventPublisher) Publish(ctx context.Context, event *Event) error {
	evJson, err := json.Marshal(event)
	if err != nil {
		return err
	}
	err = e.pubsub.Publish(ctx, e.topicName, evJson)
	return err

}

// NewEventPublisher is a function to create new EventPublisher
func NewEventPublisher(pubsub fazzpubsub.PubSub, topicName string) EventPublisher {
	return &eventPublisher{
		pubsub:    pubsub,
		topicName: topicName,
	}
}
