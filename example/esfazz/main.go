package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
	configdb "github.com/payfazz/go-apt/config"
	"github.com/payfazz/go-apt/example/esfazz/account/command"
	"github.com/payfazz/go-apt/example/esfazz/account/command/event"
	"github.com/payfazz/go-apt/example/esfazz/account/query"
	configexample "github.com/payfazz/go-apt/example/esfazz/config"
	"github.com/payfazz/go-apt/example/esfazz/migration"
	"github.com/payfazz/go-apt/pkg/esfazz/esrepo"
	"github.com/payfazz/go-apt/pkg/esfazz/eventstore/eventmongo"
	"github.com/payfazz/go-apt/pkg/esfazz/snapstore/snapmongo"
	"github.com/payfazz/go-apt/pkg/fazzdb"
	"github.com/payfazz/go-apt/pkg/fazzpubsub"
	"log"
	"sync"
)

func main() {
	fazzdb.Migrate(
		configdb.GetDB(),
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

	var wg sync.WaitGroup
	wg.Add(2)

	go queryServer(&wg, rc)
	go commandServer(&wg, rc)

	wg.Wait()
	_ = rc.Close()
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
			fmt.Printf("Message Get: %s\n", msg.Data)

			acc := &event.Account{}
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
	wg.Done()

	_ = subscription.Unsubscribe()
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
	account, err := cmd.Create(ctx, command.CreatePayload{
		Name:    "Test Account 1",
		Balance: 150,
	})
	if err != nil {
		log.Fatal(err)
	}
	sendUpdate(ctx, pubsub, account)

	// Create second account
	account, err = cmd.Create(ctx, command.CreatePayload{
		Name:    "Test Account 2",
		Balance: 100,
	})
	if err != nil {
		log.Fatal(err)
	}
	sendUpdate(ctx, pubsub, account)

	// Change name
	account, err = cmd.ChangeName(ctx, command.ChangeNamePayload{
		AccountId: account.Id,
		Name:      "New Test Account",
	})
	if err != nil {
		log.Fatal(err)
	}
	sendUpdate(ctx, pubsub, account)

	// Deposit
	account, err = cmd.Deposit(ctx, command.DepositPayload{
		AccountId: account.Id,
		Amount:    100,
	})
	if err != nil {
		log.Fatal(err)
	}
	sendUpdate(ctx, pubsub, account)

	// Withdraw
	account, err = cmd.Withdraw(ctx, command.WithdrawPayload{
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
	db := configexample.GetMongoClient().Database("command")

	eventCollection := db.Collection("events")
	_ = eventmongo.CreateAggregateUniqueIndex(eventCollection)
	snapCollection := db.Collection("snapshots")
	_ = snapmongo.CreateIdUniqueIndex(snapCollection)

	repoConfig := (&esrepo.RepositoryConfig{}).
		SetAggregateFactory(event.NewAccount).
		SetMongoEventStore(eventCollection).
		SetMongoSnapshotStore(snapCollection)
	accountRepo := esrepo.NewRepository(repoConfig)

	return command.NewAccountCommand(accountRepo)
}

func sendUpdate(ctx context.Context, pubsub fazzpubsub.PubSub, account *event.Account) {
	accountJson, _ := json.Marshal(account)
	fmt.Printf("Message Publish: %s\n", accountJson)
	_ = pubsub.Publish(ctx, "account.update", accountJson)
}

func buildContext() context.Context {
	queryDb := fazzdb.QueryDb(configdb.GetDB(),
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
