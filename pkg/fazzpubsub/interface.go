package fazzpubsub

import "context"

type Msg struct {
	Data []byte
}

type MsgHandler func(msg *Msg)

type Subscription interface {
	Unsubscribe() error
}

type PubSub interface {
	Publish(ctx context.Context, topic string, data []byte) error
	Subscribe(ctx context.Context, topic string, cb MsgHandler) (Subscription, error)
}
