package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
	"github.com/payfazz/go-apt/config"
	"github.com/payfazz/go-apt/example/esfazz/domain/account/command"
	"github.com/payfazz/go-apt/example/esfazz/domain/account/command/aggregate"
	"github.com/payfazz/go-apt/example/esfazz/domain/account/command/data"
	"github.com/payfazz/go-apt/example/esfazz/domain/account/query"
	"github.com/payfazz/go-apt/example/esfazz/migration"
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

	rc := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rc.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	go Query(&wg, rc)
	go Command(&wg, rc)

	wg.Wait()
}

func Query(wg *sync.WaitGroup, rc *redis.Client) {
	ctx := BuildContext()
	qr := query.NewAccountQuery()
	pubsub := fazzpubsub.NewRedisPubSub(rc)

	var msgWg sync.WaitGroup
	msgWg.Add(6)

	subscription, _ := pubsub.Subscribe(
		ctx, "subsname", "account.update",
		func(msg *fazzpubsub.Msg) error {
			fmt.Printf("Message: %s\n", msg.Data)

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
	defer subscription.Unsubscribe()

	msgWg.Wait()
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
	fmt.Print("]\n\n")
}

func Command(wg *sync.WaitGroup, rc *redis.Client) {
	ctx := BuildContext()
	cmd := command.NewAccountCommand()

	pubsub := fazzpubsub.NewRedisPubSub(rc)

	// Create account
	account, _ := cmd.Create(ctx, data.CreatePayload{
		Name:    "Test Account 1",
		Balance: 150,
	})
	sendUpdate(ctx, pubsub, account)

	// Create second account
	account, _ = cmd.Create(ctx, data.CreatePayload{
		Name:    "Test Account 2",
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
