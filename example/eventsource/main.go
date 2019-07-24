package main

import (
	"context"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"github.com/payfazz/go-apt/config"
	"github.com/payfazz/go-apt/example/eventsource/domain/account/command"
	"github.com/payfazz/go-apt/example/eventsource/domain/account/command/aggregate"
	"github.com/payfazz/go-apt/example/eventsource/domain/account/command/data"
	"github.com/payfazz/go-apt/example/eventsource/domain/account/query"
	"github.com/payfazz/go-apt/example/eventsource/migration"
	"github.com/payfazz/go-apt/pkg/fazzdb"
	"github.com/payfazz/go-apt/pkg/fazzpubsub"
	"sync"
)

func main() {
	fazzdb.Migrate(
		config.GetDB(),
		"es-example",
		true,
		true,
		migration.Version1,
	)

	var wg sync.WaitGroup
	wg.Add(2)

	go Query(&wg)
	go Command(&wg)

	wg.Wait()
}

func Command(wg *sync.WaitGroup) {
	ctx := BuildContext()
	cmd := command.NewAccountCommand()

	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Close()
	pubsub := fazzpubsub.NewNatsPubSub(nc)

	// Create account
	account, _ := cmd.Create(ctx, data.CreatePayload{
		Name:    "Test Account",
		Balance: 100,
	})
	sendUpdate(ctx, pubsub, account)

	// Change name
	account, _ = cmd.ChangeName(ctx, data.ChangeNamePayload{
		AccountId: account.Id,
		Name:      "New Test Account",
	})
	sendUpdate(ctx, pubsub, account)

	// Deposit
	account, _ = cmd.Deposit(ctx, data.DepositPayload{
		AccountId: account.Id,
		Amount:    100,
	})
	sendUpdate(ctx, pubsub, account)

	// Withdraw
	account, _ = cmd.Withdraw(ctx, data.WithdrawPayload{
		AccountId: account.Id,
		Amount:    200,
	})
	sendUpdate(ctx, pubsub, account)

	// Delete
	account, _ = cmd.Delete(ctx, account.Id)
	sendUpdate(ctx, pubsub, account)

	wg.Done()
}

func sendUpdate(ctx context.Context, pubsub fazzpubsub.PubSub, account *aggregate.Account) {
	accountJson, _ := json.Marshal(account)
	_ = pubsub.Publish(ctx, "account.update", accountJson)
}

func Query(wg *sync.WaitGroup) {
	qr := query.NewAccountQuery()

	nc, _ := nats.Connect(nats.DefaultURL)
	pubsub := fazzpubsub.NewNatsPubSub(nc)
	var msgWg sync.WaitGroup
	msgWg.Add(5)

	subscription, _ := pubsub.Subscribe(
		context.Background(),
		"subsname",
		"account.update",
		func(msg *fazzpubsub.Msg) error {
			fmt.Printf("Message: %s\n", msg.Data)
			ctx := BuildContext()

			acc := &aggregate.Account{}
			err := json.Unmarshal(msg.Data, acc)
			if err != nil {
				return err
			}

			err = qr.DirectUpdate(ctx, acc)
			if err != nil {
				return err
			}

			printState(ctx, qr)
			msgWg.Done()
			return nil
		},
	)

	msgWg.Wait()
	_ = subscription.Unsubscribe()
	nc.Close()
	wg.Done()
}

func printState(ctx context.Context, qr query.AccountQuery) {
	accList, err := qr.All(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Current Account List: [")
	for _, acc := range accList {
		fmt.Printf("(%s, %d),", acc.Name, acc.Balance)
	}
	fmt.Println("]")
}

func BuildContext() context.Context {
	queryDb := fazzdb.QueryDb(config.GetDB(),
		fazzdb.Config{
			Limit:           20,
			Offset:          0,
			Lock:            fazzdb.LO_NONE,
			DevelopmentMode: true,
		})

	ctx := context.Background()
	ctx = fazzdb.NewQueryContext(ctx, queryDb)

	return ctx
}
