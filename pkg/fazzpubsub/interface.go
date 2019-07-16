package fazzpubsub

import "context"

type Msg struct {
	Topic string
	Data  []byte
}

type MsgHandler func(msg *Msg) error

type Subscription interface {
	Unsubscribe() error
}

type PubSub interface {
	Publish(ctx context.Context, topic string, data []byte) error
	Subscribe(ctx context.Context, name string, topic string, cb MsgHandler) (Subscription, error)
}
