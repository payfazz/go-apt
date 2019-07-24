package event

const (
	ACCOUNT_CREATED      = "account.created"
	ACCOUNT_NAME_CHANGED = "account.name.changed"
	ACCOUNT_DEPOSITED    = "account.deposited"
	ACCOUNT_WITHDRAWN    = "account.withdrawn"
	ACCOUNT_DELETED      = "account.deleted"
)

type AccountCreatedData struct {
	Name    string `json:"name"`
	Balance int64  `json:"balance"`
}

type AccountNameChangedData struct {
	Name string `json:"text"`
}

type AccountDepositedData struct {
	Amount int64 `json:"amount"`
}

type AccountWithdrawnData struct {
	Amount int64 `json:"amount"`
}
