package main

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/payfazz/go-apt/pkg/fazzpubsub"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	ctx := context.Background()
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Close()
	pubsub := fazzpubsub.NewNatsPubSub(nc)

	wg.Add(3)
	subs, _ := pubsub.Subscribe(ctx, "subs", "event", func(msg *fazzpubsub.Msg) error {
		fmt.Printf("%s\n", msg.Data)
		wg.Done()
		return nil
	})
	_ = pubsub.Publish(ctx, "event", []byte("test 1"))
	_ = pubsub.Publish(ctx, "event", []byte("test 2"))
	_ = pubsub.Publish(ctx, "event", []byte("test 3"))
	wg.Wait()
	_ = subs.Unsubscribe()
}
