package model

import "time"

// User represents a user in the system
type User struct {
	ID        uint64    `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
