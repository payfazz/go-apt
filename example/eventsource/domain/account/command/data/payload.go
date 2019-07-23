package data

type CreatePayload struct {
	Name    string `json:"name"`
	Balance int64  `json:"balance"`
}

type ChangeNamePayload struct {
	AccountId string `json:"account_id"`
	Name      string `json:"name"`
}

type DepositPayload struct {
	AccountId string `json:"account_id"`
	Amount    int64  `json:"amount"`
}

type WithdrawPayload struct {
	AccountId string `json:"account_id"`
	Amount    int64  `json:"amount"`
}
