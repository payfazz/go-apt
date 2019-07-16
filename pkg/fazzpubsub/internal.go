package fazzpubsub

import (
	"context"
)

type internalPubSub struct {
	subs map[string]map[string]subscription
}

type subscription struct {
	pubsub internalPubSub
	name   string
	topic  string
	cb     MsgHandler
}

func (p *internalPubSub) Publish(ctx context.Context, topic string, data []byte) error {
	for _, subs := range p.subs[topic] {
		_ = subs.cb(&Msg{
			Topic: topic,
			Data:  data,
		})
	}
	return nil
}

func (p *internalPubSub) Subscribe(ctx context.Context, name string, topic string, cb MsgHandler) (Subscription, error) {
	subs := subscription{
		name:  name,
		topic: topic,
		cb:    cb,
	}
	subsTopic, exists := p.subs[topic]
	if !exists {
		subsTopic = make(map[string]subscription)
		p.subs[topic] = subsTopic
	}
	subsTopic[name] = subs
	return &subs, nil
}

func (s *subscription) Unsubscribe() error {
	delete(s.pubsub.subs[s.topic], s.name)
	return nil
}

// NewInternalPubSub create internal pub sub, use only for test
func NewInternalPubSub() PubSub {
	return &internalPubSub{
		subs: make(map[string]map[string]subscription),
	}
}
