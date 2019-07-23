package main

import (
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/payfazz/go-apt/config"
	"github.com/payfazz/go-apt/example/eventsource/domain/account/command"
	"github.com/payfazz/go-apt/example/eventsource/domain/account/command/data"
	"github.com/payfazz/go-apt/example/eventsource/domain/account/query"
	"github.com/payfazz/go-apt/example/eventsource/migration"
	"github.com/payfazz/go-apt/pkg/fazzdb"
)

func main() {
	fazzdb.Migrate(
		config.GetDB(),
		"es-example",
		true,
		true,
		migration.Version1,
	)

	ctx := BuildContext()
	cmd := command.NewAccountCommand()
	qr := query.NewAccountQuery()

	accountId, _ := cmd.Create(ctx, data.CreatePayload{
		Name:    "Test Account",
		Balance: 100,
	})
	SendUpdateToQuery(ctx, cmd, qr, *accountId)

	accountModel, _ := qr.Get(ctx, *accountId)
	fmt.Printf("Name: %s, Balance: %d\n", accountModel.Name, accountModel.Balance)

	_ = cmd.ChangeName(ctx, data.ChangeNamePayload{
		AccountId: *accountId,
		Name:      "New Test Account",
	})
	SendUpdateToQuery(ctx, cmd, qr, *accountId)

	accountModel, _ = qr.Get(ctx, *accountId)
	fmt.Printf("Name: %s, Balance: %d\n", accountModel.Name, accountModel.Balance)

	_ = cmd.Deposit(ctx, data.DepositPayload{
		AccountId: *accountId,
		Amount:    100,
	})
	SendUpdateToQuery(ctx, cmd, qr, *accountId)

	accountModel, _ = qr.Get(ctx, *accountId)
	fmt.Printf("Name: %s, Balance: %d\n", accountModel.Name, accountModel.Balance)

	_ = cmd.Withdraw(ctx, data.WithdrawPayload{
		AccountId: *accountId,
		Amount:    200,
	})
	SendUpdateToQuery(ctx, cmd, qr, *accountId)

	accountModel, _ = qr.Get(ctx, *accountId)
	fmt.Printf("Name: %s, Balance: %d\n", accountModel.Name, accountModel.Balance)

	_ = cmd.Delete(ctx, *accountId)
	SendUpdateToQuery(ctx, cmd, qr, *accountId)

	accountList, _ := qr.All(ctx)
	fmt.Printf("Account Count: %d\n", len(accountList))
}

func SendUpdateToQuery(ctx context.Context, cmd command.AccountCommand, qr query.AccountQuery, id string) {
	account, _ := cmd.DirectGet(ctx, id)
	_ = qr.DirectUpdate(ctx, account)
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
