package swagger

// Common
type UserNotFoundError struct {
	Success bool   `json:"success" example:"false"`
	Code    string `json:"code" example:"404"`
	Message string `json:"message" example:"User not found"`
}

type InvalidUserIDError struct {
	Success bool   `json:"success" example:"false"`
	Code    string `json:"code" example:"400"`
	Message string `json:"message" example:"Invalid user ID"`
}

// Balance Specific
type BalanceErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Code    string `json:"code" example:"500"`
	Message string `json:"message" example:"Failed to retrieve balance"`
}

// Transaction Specific
type InvalidAmountError struct {
	Success bool   `json:"success" example:"false"`
	Code    string `json:"code" example:"400"`
	Message string `json:"message" example:"Invalid amount format"`
}

type InvalidTransactionStateError struct {
	Success bool   `json:"success" example:"false"`
	Code    string `json:"code" example:"400"`
	Message string `json:"message" example:"Invalid transaction state"`
}

type InsufficientBalanceError struct {
	Success bool   `json:"success" example:"false"`
	Code    string `json:"code" example:"422"`
	Message string `json:"message" example:"Insufficient balance"`
}

type InvalidSourceTypeError struct {
	Success bool   `json:"success" example:"false"`
	Code    string `json:"code" example:"400"`
	Message string `json:"message" example:"Invalid Source-Type header"`
}

type MissingSourceTypeError struct {
	Success bool   `json:"success" example:"false"`
	Code    string `json:"code" example:"404"`
	Message string `json:"message" example:"Source-Type header not found"`
}

type TransactionResponse struct {
	Success     bool   `json:"success" example:"true"`
	UserID      uint64 `json:"userId" example:"42"`
	Balance     string `json:"balance" example:"125.75"`
	Transaction string `json:"transactionId" example:"tx_abc123def456"`
	ProcessedAt string `json:"processedAt" example:"2025-03-14T14:30:45Z"`
}
