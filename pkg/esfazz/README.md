# esfazz

## Usage example

#### Implement Aggregate for aggregation model
```go
type Account struct {
  esfazz.BaseAggregate
  Name      string     `json:"name"`
  Balance   int64      `json:"balance"`
}

// Apply function implement Aggregate inteface
func (a *Account) Apply(evs ...*esfazz.Event) error {
  for _, ev := range evs {
    a.Version = a.Version + 1
    
    // apply base on event type
    switch ev.Type {
    
    case "account.created":
      var data *AccountCreatedData
      // .... unmarshal data from ev.Data
      a.Id = ev.Aggregate.GetId()
      a.Name = data.Name
      a.Balance = data.Balance
      
    case "account.withdrawn":
      var data *AccountWithdrawnData
      // .... unmarshal data from ev.Data
      a.Balance -= data.Amount
    }
  }
  return nil
}

// NewAccount used for aggregate constructor
func NewAccount(id string) esfazz.Aggregate {
  account := &Account{}
  account.Id = id
  return account
}
```

#### Setup repository

- MongoDB
```go
// get mongo db collections and create index
db := mongoclient.Database("dbname")

eventCollection := db.Collection("events")
eventmongo.CreateAggregateUniqueIndex(eventCollection)

snapCollection := db.Collection("snapshots")
snapmongo.CreateIdUniqueIndex(snapCollection)

// create repository
repoConfig := (&esrepo.RepositoryConfig{}).
  SetAggregateFactory(NewAccount).
  SetMongoEventStore(eventCollection).
  SetMongoSnapshotStore(snapCollection)
accountRepo := esrepo.NewRepository(repoConfig)
```
- PostgreSQL
```go
// setup database migration
migrationVersion := fazzdb.MigrationVersion{
  Tables: []*fazzdb.MigrationTable{
    eventpostgres.CreateEventsTable("events"),
    snappostgres.CreateSnapshotsTable("snapshots"),
  },
}
// .... run migration to database

// create repository
repoConfig := (&esrepo.RepositoryConfig{}).
  SetAggregateFactory(NewAccount).
  SetPostgresEventStore("events").
  SetPostgresSnapshotStore("snapshots")
accountRepo := esrepo.NewRepository(repoConfig)
```

#### Save event
```go
func CreateAccount(id string, name string, balance int){
  ev := &esfazz.EventPayload{
    Type: "account.created",
    Aggregate: NewAccount(id),
    Data: AccountCreatedData{
      Name:      name,
      Balance:   balance,
    },
  }
  err := accountRepo.Save(ctx, ev)
}
```

#### Find Aggregate And Then Save New Event

```go
func WithdrawAccount(id string, amount int) {
  agg, err := a.repository.Find(ctx, id)
  if err != nil { 
    // handle error
  }
  account := agg.(*Account)
  if account == nil {
    // handle null
  }

  // validate aggregate before creating new event
  if account.Balance < amount {
    // handle error
  }

  // use aggregate for saving newer event
  ev := &esfazz.EventPayload{
    Type: "account.withdrawn",
    Aggregate: account,
    Data: AccountWithdrawnData{Amount: amount},
  }
  err := accountRepo.Save(ctx, ev)
}
```
