package repository

import (
	"gorm.io/gorm"
)

// RepositoryManager manages all repositories
type RepositoryManager struct {
	User        UserInterface
	Balance     BalanceInterface
	Transaction TransactionInterface
}

// NewRepositoryManager creates a new repository manager
func New(db *gorm.DB) *RepositoryManager {
	return &RepositoryManager{
		User:        NewUserRepository(db),
		Balance:     NewBalanceRepository(db),
		Transaction: NewTransactionRepository(db),
	}
}
