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
	repocommand "github.com/payfazz/go-apt/example/esfazz/domain/account/command/repository"
	"github.com/payfazz/go-apt/example/esfazz/domain/account/query"
	"github.com/payfazz/go-apt/example/esfazz/migration"
	"github.com/payfazz/go-apt/pkg/esfazz/eventstore/eventmongo"
	"github.com/payfazz/go-apt/pkg/esfazz/snapstore/snapmongo"
	"github.com/payfazz/go-apt/pkg/fazzdb"
	"github.com/payfazz/go-apt/pkg/fazzpubsub"
	"log"
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
		Password: "cashfazz",
		DB:       0,
	})
	defer rc.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	go queryServer(&wg, rc)
	go commandServer(&wg, rc)

	wg.Wait()
}

func queryServer(wg *sync.WaitGroup, rc *redis.Client) {
	ctx := buildContext()
	qr := query.NewAccountQuery()
	pubsub := fazzpubsub.RedisPubSub(rc)

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

func commandServer(wg *sync.WaitGroup, rc *redis.Client) {
	ctx := buildContext()
	cmd := createCommand()
	pubsub := fazzpubsub.RedisPubSub(rc)

	// Create account
	account, err := cmd.Create(ctx, data.CreatePayload{
		Name:    "Test Account 1",
		Balance: 150,
	})
	if err != nil {
		log.Fatal(err)
	}
	sendUpdate(ctx, pubsub, account)

	// Create second account
	account, err = cmd.Create(ctx, data.CreatePayload{
		Name:    "Test Account 2",
		Balance: 100,
	})
	if err != nil {
		log.Fatal(err)
	}
	sendUpdate(ctx, pubsub, account)

	// Change name
	account, err = cmd.ChangeName(ctx, data.ChangeNamePayload{
		AccountId: account.Id,
		Name:      "New Test Account",
	})
	if err != nil {
		log.Fatal(err)
	}
	sendUpdate(ctx, pubsub, account)

	// Deposit
	account, err = cmd.Deposit(ctx, data.DepositPayload{
		AccountId: account.Id,
		Amount:    100,
	})
	if err != nil {
		log.Fatal(err)
	}
	sendUpdate(ctx, pubsub, account)

	// Withdraw
	account, err = cmd.Withdraw(ctx, data.WithdrawPayload{
		AccountId: account.Id,
		Amount:    200,
	})
	if err != nil {
		log.Fatal(err)
	}
	sendUpdate(ctx, pubsub, account)

	// Delete
	account, err = cmd.Delete(ctx, account.Id)
	if err != nil {
		log.Fatal(err)
	}
	sendUpdate(ctx, pubsub, account)

	wg.Done()
}

func createCommand() command.AccountCommand {
	db := GetMongoClient().Database("command")

	eventCollection := db.Collection("events")
	_ = eventmongo.CreateAggregateUniqueIndex(eventCollection)
	snapCollection := db.Collection("snapshots")
	_ = snapmongo.CreateIdUniqueIndex(snapCollection)

	eventStore := eventmongo.EventStore(eventCollection)
	snapStore := snapmongo.SnapshotStore(snapCollection)

	repo := repocommand.NewAccountEventRepository(eventStore, snapStore)

	return command.NewAccountCommand(repo)
}

func sendUpdate(ctx context.Context, pubsub fazzpubsub.PubSub, account *aggregate.Account) {
	accountJson, _ := json.Marshal(account)
	_ = pubsub.Publish(ctx, "account.update", accountJson)
}

func buildContext() context.Context {
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
