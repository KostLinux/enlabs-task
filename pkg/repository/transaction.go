package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"enlabs-task/pkg/model"
)

// TransactionRepository handles operations on transaction data
type TransactionRepository struct {
	db *sqlx.DB
}

// NewTransactionRepository creates a new TransactionRepository instance
func NewTransactionRepository(db *sqlx.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Create saves a new transaction record
func (r *TransactionRepository) Create(tx *model.Transaction) error {
	query := `
        INSERT INTO transactions 
        (transaction_id, user_id, state, amount, source_type, previous_balance, new_balance)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at
    `

	row := r.db.QueryRow(
		query,
		tx.TransactionID,
		tx.UserID,
		tx.State,
		tx.Amount,
		tx.SourceType,
		tx.PreviousBalance,
		tx.NewBalance,
	)

	if err := row.Scan(&tx.ID, &tx.CreatedAt); err != nil {
		return fmt.Errorf("error creating transaction: %w", err)
	}

	return nil
}

// FindByTransactionID checks if a transaction with the given ID exists
func (r *TransactionRepository) FindByTransactionID(transactionID string) (*model.Transaction, error) {
	var tx model.Transaction

	query := `SELECT * FROM transactions WHERE transaction_id = $1`
	err := r.db.Get(&tx, query, transactionID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Transaction doesn't exist
		}
		return nil, fmt.Errorf("error fetching transaction: %w", err)
	}

	return &tx, nil
}
