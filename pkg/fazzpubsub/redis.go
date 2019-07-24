package fazzpubsub

import (
	"context"
	"github.com/go-redis/redis"
)

type redisPubSub struct {
	rc *redis.Client
}

type redisSubscription struct {
	p *redis.PubSub
}

func (r *redisPubSub) Publish(ctx context.Context, topic string, data []byte) error {
	_, err := r.rc.Publish(topic, data).Result()
	return err
}

func (r *redisPubSub) Subscribe(ctx context.Context, name string, topic string, cb MsgHandler) (Subscription, error) {
	p := r.rc.Subscribe(topic)
	msgChan := p.Channel()
	go func() {
		for redMsg := range msgChan {
			msg := &Msg{
				Topic: redMsg.Channel,
				Data:  []byte(redMsg.Payload),
			}
			_ = cb(msg)
		}
	}()
	return &redisSubscription{p: p}, nil
}

func (r *redisSubscription) Unsubscribe() error {
	err := r.p.Close()
	return err
}

func NewRedisPubSub(rc *redis.Client) PubSub {
	return &redisPubSub{rc: rc}
}
