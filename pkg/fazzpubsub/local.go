package fazzpubsub

import (
	"context"
)

type localPubSub struct {
	subs map[string]map[string]subscription
}

type subscription struct {
	pubsub localPubSub
	name   string
	topic  string
	cb     MsgHandler
}

func (p *localPubSub) Publish(ctx context.Context, topic string, data []byte) error {
	for _, subs := range p.subs[topic] {
		_ = subs.cb(&Msg{
			Topic: topic,
			Data:  data,
		})
	}
	return nil
}

func (p *localPubSub) Subscribe(ctx context.Context, name string, topic string, cb MsgHandler) (Subscription, error) {
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

// LocalPubSub create internal pub sub, use only for test
func LocalPubSub() PubSub {
	return &localPubSub{
		subs: make(map[string]map[string]subscription),
	}
}
