package repository

import (
	"enlabs-task/pkg/model"

	"gorm.io/gorm"
)

// UserRepositoryInterface defines the operations for User entities
type UserInterface interface {
	GetByID(id uint64) (*model.User, error)
	Exists(id uint64) (bool, error)
}

// BalanceRepositoryInterface defines the operations for Balance entities
type BalanceInterface interface {
	GetByUserID(userID uint64) (*model.Balance, error)
	UpdateAmount(userID uint64, newAmount float64) error
}

// TransactionRepositoryInterface defines the operations for Transaction entities
type TransactionInterface interface {
	Create(tx *model.Transaction) error
	FindByTransactionID(transactionID string) (*model.Transaction, error)
}

// RepositoryManager manages all repositories
type RepositoryManager struct {
	User        *UserRepository
	Balance     *BalanceRepository
	Transaction *TransactionRepository
}

// NewRepositoryManager creates a new repository manager
func NewRepositoryManager(db *gorm.DB) *RepositoryManager {
	return &RepositoryManager{
		User:        NewUserRepository(db),
		Balance:     NewBalanceRepository(db),
		Transaction: NewTransactionRepository(db),
	}
}
