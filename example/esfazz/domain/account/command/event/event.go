package event

const (
	// AccountCreatedType is string type for account created
	AccountCreatedType = "account.created"
	// AccountNameChangedType is string type for account name changed
	AccountNameChangedType = "account.name.changed"
	// AccountDepositedType is string type for account deposited
	AccountDepositedType = "account.deposited"
	// AccountWithdrawnType is string type for account withdrawn
	AccountWithdrawnType = "account.withdrawn"
	// AccountDeletedType is string type for account deleted
	AccountDeletedType = "account.deleted"
)

// AccountCreatedData is account event created data
type AccountCreatedData struct {
	Name    string `json:"name"`
	Balance int64  `json:"balance"`
}

// AccountNameChangedData is account name changed event data
type AccountNameChangedData struct {
	Name string `json:"text"`
}

// AccountDepositedData is account deposited event data
type AccountDepositedData struct {
	Amount int64 `json:"amount"`
}

// AccountWithdrawnData is account withdrawn event data
type AccountWithdrawnData struct {
	Amount int64 `json:"amount"`
}
