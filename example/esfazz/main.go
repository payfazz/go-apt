package main

import (
	"context"
	"fmt"
	_ "github.com/lib/pq"
	configdb "github.com/payfazz/go-apt/config"
	"github.com/payfazz/go-apt/example/esfazz/account/command"
	"github.com/payfazz/go-apt/example/esfazz/account/model"
	"github.com/payfazz/go-apt/example/esfazz/account/query"
	"github.com/payfazz/go-apt/example/esfazz/migration"
	"github.com/payfazz/go-apt/pkg/esfazz/esrepo"
	"github.com/payfazz/go-apt/pkg/esfazz/eventstore/eventpostgres"
	"github.com/payfazz/go-apt/pkg/esfazz/snapstore/snappostgres"
	"github.com/payfazz/go-apt/pkg/fazzdb"
	"log"
	"strings"
)

func main() {
	fazzdb.Migrate(
		configdb.GetDB(),
		"es-example",
		true,
		true,
		migration.Version1,
	)

	ctx := buildContext()
	accountCommand := createCommand()
	accountQuery := query.NewAccountQuery()

	// Create account
	_, err := accountCommand.Create(ctx, command.CreatePayload{
		Name:    "Test Account 1",
		Balance: 150,
	})
	if err != nil {
		log.Fatal(err)
	}
	printState(ctx, accountQuery)

	// Create second account
	id, err := accountCommand.Create(ctx, command.CreatePayload{
		Name:    "Test Account 2",
		Balance: 100,
	})
	if err != nil {
		log.Fatal(err)
	}
	printState(ctx, accountQuery)

	// Change name
	err = accountCommand.ChangeName(ctx, command.ChangeNamePayload{
		AccountId: *id,
		Name:      "Build Test Account",
	})
	if err != nil {
		log.Fatal(err)
	}
	printState(ctx, accountQuery)

	// Deposit
	err = accountCommand.Deposit(ctx, command.DepositPayload{
		AccountId: *id,
		Amount:    100,
	})
	if err != nil {
		log.Fatal(err)
	}
	printState(ctx, accountQuery)

	// Withdraw
	err = accountCommand.Withdraw(ctx, command.WithdrawPayload{
		AccountId: *id,
		Amount:    200,
	})
	if err != nil {
		log.Fatal(err)
	}
	printState(ctx, accountQuery)

	// Delete
	err = accountCommand.Delete(ctx, *id)
	if err != nil {
		log.Fatal(err)
	}
	printState(ctx, accountQuery)

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

func createCommand() command.AccountCommand {
	//db := configexample.GetMongoClient().Database("command")

	//eventCollection := db.Collection("events")
	//_ = eventmongo.CreateAggregateUniqueIndex(eventCollection)
	//snapCollection := db.Collection("snapshots")
	//_ = snapmongo.CreateIdUniqueIndex(snapCollection)

	eventStore := eventpostgres.EventStore("events")
	//eventStore := eventmongo.EventStore(eventCollection)

	snapshotStore := snappostgres.ModelSnapshotStore(model.NewAccount)
	//snapshotStore := snappostgres.JSONSnapshotStore("snapshots")
	//snapshotStore := snapmongo.SnapshotStore(snapCollection)

	repoConfig := esrepo.Config{}.
		SetAggregateFactory(model.NewAccount).
		SetEventStore(eventStore).
		SetSnapshotStore(snapshotStore).
		AddEventListener(esrepo.SnapshotSaver(snapshotStore, model.NewAccount))
	accountRepo := esrepo.Build(repoConfig)

	return command.NewAccountCommand(accountRepo)
}

func printState(ctx context.Context, accountQuery query.AccountQuery) {
	accounts, err := accountQuery.All(ctx)
	if err != nil {
		log.Println(err)
	}

	var sb strings.Builder
	for idx, account := range accounts {
		if idx == 0 {
			sb.WriteString("Current Account List: [ ")
		} else {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("(%s, %d)", account.Name, account.Balance))
	}
	sb.WriteString(" ]")

	log.Println(sb.String())
}
