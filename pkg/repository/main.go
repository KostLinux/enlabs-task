package repository

import (
	"github.com/jmoiron/sqlx"

	"enlabs-task/pkg/model"
)

// UserRepositoryInterface defines the operations for User entities
type UserRepositoryInterface interface {
	GetByID(id uint64) (*model.User, error)
	Exists(id uint64) (bool, error)
}

// BalanceRepositoryInterface defines the operations for Balance entities
type BalanceRepositoryInterface interface {
	GetByUserID(userID uint64) (*model.Balance, error)
	UpdateAmount(userID uint64, newAmount float64) error
}

// TransactionRepositoryInterface defines the operations for Transaction entities
type TransactionRepositoryInterface interface {
	Create(tx *model.Transaction) error
	FindByTransactionID(transactionID string) (*model.Transaction, error)
}

// RepositoryManager holds all repository instances
type RepositoryManager struct {
	User        UserRepositoryInterface
	Balance     BalanceRepositoryInterface
	Transaction TransactionRepositoryInterface
}

// NewRepositoryManager creates a new RepositoryManager with initialized repositories
func NewRepositoryManager(db *sqlx.DB) *RepositoryManager {
	return &RepositoryManager{
		User:        NewUserRepository(db),
		Balance:     NewBalanceRepository(db),
		Transaction: NewTransactionRepository(db),
	}
}
