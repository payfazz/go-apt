package command

// CreatePayload is create account function payload
type CreatePayload struct {
	Name    string `json:"name"`
	Balance int64  `json:"balance"`
}

// ChangeNamePayload is change name function payload
type ChangeNamePayload struct {
	AccountId string `json:"account_id"`
	Name      string `json:"name"`
}

// DepositPayload is deposit account function payload
type DepositPayload struct {
	AccountId string `json:"account_id"`
	Amount    int64  `json:"amount"`
}

// WithdrawPayload is withdraw account function payload
type WithdrawPayload struct {
	AccountId string `json:"account_id"`
	Amount    int64  `json:"amount"`
}
