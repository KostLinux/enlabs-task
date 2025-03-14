package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"enlabs-task/pkg/model"
)

// Transaction defines the transaction repository interface
type TransactionInterface interface {
	Create(tx *model.Transaction) error
	FindByTransactionID(transactionID string) (*model.Transaction, error)
}

// TransactionRepository handles operations on transaction data
type TransactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository creates a new TransactionRepository instance
func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Create saves a new transaction record
func (repository *TransactionRepository) Create(tx *model.Transaction) error {
	result := repository.db.Create(tx)
	if result.Error != nil {
		return fmt.Errorf("error creating transaction: %w", result.Error)
	}

	return nil
}

// FindByTransactionID checks if a transaction with the given ID exists
func (repository *TransactionRepository) FindByTransactionID(transactionID string) (*model.Transaction, error) {
	var tx model.Transaction

	result := repository.db.Where("transaction_id = ?", transactionID).First(&tx)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Transaction doesn't exist - this is an expected case, not an error
			return nil, nil
		}
		// Unexpected database error
		return nil, fmt.Errorf("error fetching transaction: %w", result.Error)
	}

	return &tx, nil
}
