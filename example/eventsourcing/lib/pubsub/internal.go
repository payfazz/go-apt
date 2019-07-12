package pubsub

import (
	"context"
	"sync"
)

type Msg struct {
	Data []byte
}

type MsgHandler func(msg *Msg)

type Subscription interface {
	Unsubscribe() error
}

type PubSub interface {
	Publish(ctx context.Context, subject string, data []byte) error
	Subscribe(ctx context.Context, subject string, cb MsgHandler) (Subscription, error)
}

type internalPubSub struct {
	subs map[string][]subscription
}

type subscription struct {
	pubsub  internalPubSub
	subject string
	cb      MsgHandler
}

func (p *internalPubSub) Publish(ctx context.Context, subject string, data []byte) error {
	for _, subs := range p.subs[subject] {
		subs.cb(&Msg{Data: data})
	}
	return nil
}

func (p *internalPubSub) Subscribe(ctx context.Context, subject string, cb MsgHandler) (Subscription, error) {
	subs := subscription{
		subject: subject,
		cb:      cb,
	}
	p.subs[subject] = append(p.subs[subject], subs)
	return &subs, nil
}

func (s *subscription) Unsubscribe() error {
	return nil
}

func NewInternalPubSub() PubSub {
	return &internalPubSub{
		subs: map[string][]subscription{},
	}
}

var ps PubSub
var psOnce sync.Once

func SingletonPubSub() PubSub {
	psOnce.Do(func() {
		ps = NewInternalPubSub()
	})
	return ps
}
