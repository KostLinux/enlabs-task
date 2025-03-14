package model

import (
	"time"
)

// Transaction represents a financial transaction
type Transaction struct {
	ID              uint64    `json:"id" db:"id"`
	TransactionID   string    `json:"transaction_id" db:"transaction_id"`
	UserID          uint64    `json:"user_id" db:"user_id"`
	State           string    `json:"state" db:"state"`
	Amount          float64   `json:"amount" db:"amount"`
	SourceType      string    `json:"source_type" db:"source_type"`
	PreviousBalance float64   `json:"previous_balance" db:"previous_balance"`
	NewBalance      float64   `json:"new_balance" db:"new_balance"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// TransactionRequest represents the incoming transaction request
type TransactionRequest struct {
	State         string `json:"state" example:"win" binding:"required,oneof=win lose"`
	Amount        string `json:"amount" example:"10.50" binding:"required"`
	TransactionID string `json:"transactionId" example:"tx_abc123def456" binding:"required"`
}

// TransactionResponse is the response for the POST transaction endpoint
// Balance is string due to formatted with 2 decimal places
// Float64 is not used to avoid floating point precision during calculation
// and JSON Marshalling
type TransactionResponse struct {
	Success     bool   `json:"success" example:"true"`
	UserID      uint64 `json:"userId" example:"42"`
	Balance     string `json:"balance" example:"125.75"`
	Transaction string `json:"transactionId" example:"tx_abc123def456"`
	ProcessedAt string `json:"processedAt" example:"2025-03-14T14:30:45Z"`
}
