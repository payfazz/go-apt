# fazzeventsource

## Usage example

#### Implement Aggregate for aggregation model
```go

type Account struct {
  fazzeventsource.Entity
  Name      string
  Balance   int64
}

// Apply function apply event to domain
func (a *Account) Apply(event fazzeventsource.Event) error {
    a.Entity.Apply(event)
    switch event.Type {
    case "account.created":
      a.Id = ev.Aggregate.GetId()
      a.Name = event.Data["name"]
      a.Balance = event.Data["balance"]
      
    case "account.withdrawn":
      a.Balance -= event.Data["amount"]
    }
    return nil
}
```

#### Create repository

```go
// setup database migration
migrationVersion := fazzdb.MigrationVersion{
  Tables: []*fazzdb.MigrationTable{
    fazzeventsource.CreateEventsTable("account_events"),
  },
}
// .... run migration to database

// create repository
accountRepo := fazzeventsource.NewRepository("account_events")
```

#### Save event into entity
```go
func CreateAccount(ctx context.Context, id string, name string, balance int){
    acc := Account{}
    ev := fazzeventsource.Event{
        Type: "account.created",
        Data: map[string]interface{}{
            "name":      name,
            "balance":   balance,
        },
    }
    err := acc.apply(ev)
    // handle err
    err = accountRepo.Commit(ctx, acc)
    // handle err
}
```

### Load entity from event repo

```go
func WithdrawAccount(ctx context.Context, id string, name string, amount int){
    acc := Account{}
    acc.Id = id
    err := accountRepo.GetLatest(acc)
    // handle err
    ev := fazzeventsource.Event{
        Type: "account.withdrawn",
        Data: map[string]interface{}{
            "amount":   amount,
        },
    }
    err := acc.apply(ev)
    // handle err
    err = accountRepo.Commit(ctx, acc)
    // handle err
}
```
