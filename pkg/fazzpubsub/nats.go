package fazzpubsub

import (
	"context"
	"github.com/nats-io/nats.go"
)

type natsPubSub struct {
	nc *nats.Conn
}

func (p *natsPubSub) Publish(ctx context.Context, topic string, data []byte) error {
	return p.nc.Publish(topic, data)
}

func (p *natsPubSub) Subscribe(ctx context.Context, name string, topic string, cb MsgHandler) (Subscription, error) {
	return p.nc.Subscribe(topic, func(msg *nats.Msg) {
		_ = cb(&Msg{
			Topic: msg.Subject,
			Data:  msg.Data,
		})
	})
}

func NewNatsPubSub(nc *nats.Conn) PubSub {
	return &natsPubSub{
		nc: nc,
	}
}
