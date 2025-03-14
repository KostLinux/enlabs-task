package model

import (
	"time"
)

type Balance struct {
	ID        uint64    `json:"id" db:"id"`
	UserID    uint64    `json:"user_id" db:"user_id"`
	Amount    float64   `json:"amount" db:"amount"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// BalanceResponse is the response for the GET balance endpoint
// Balance is string due to formatted with 2 decimal places
// Float64 is not used to avoid floating point precision during calculation
// and JSON Marshalling
type BalanceResponse struct {
	UserID  uint64 `json:"userId" example:"12345"`
	Balance string `json:"balance" example:"100.00"`
}
