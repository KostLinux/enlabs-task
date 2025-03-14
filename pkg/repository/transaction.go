package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"enlabs-task/pkg/model"
)

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
			return nil, fmt.Errorf("transaction with ID %s not found", transactionID)
		}
		return nil, fmt.Errorf("error fetching transaction: %w", result.Error)
	}

	return &tx, nil
}
