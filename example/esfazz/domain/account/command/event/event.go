package event

const (
	ACCOUNT_CREATED      = "account.created"
	ACCOUNT_NAME_CHANGED = "account.name.changed"
	ACCOUNT_DEPOSITED    = "account.deposited"
	ACCOUNT_WITHDRAWN    = "account.withdrawn"
	ACCOUNT_DELETED      = "account.deleted"
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
