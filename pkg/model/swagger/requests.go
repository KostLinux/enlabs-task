package swagger

// TransactionRequest represents the transaction request payload
type TransactionRequest struct {
	State         string `json:"state" enums:"win,lose" example:"win"`
	Amount        string `json:"amount" example:"50.25"`
	TransactionID string `json:"transactionId" example:"tx_win_12345"`
}
