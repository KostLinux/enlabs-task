package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"enlabs-task/pkg/model"
)

// BalanceRepository handles operations on balance data
type BalanceRepository struct {
	db *sqlx.DB
}

// NewBalanceRepository creates a new BalanceRepository instance
func NewBalanceRepository(db *sqlx.DB) *BalanceRepository {
	return &BalanceRepository{db: db}
}

// GetByUserID retrieves a user's balance
func (r *BalanceRepository) GetByUserID(userID uint64) (*model.Balance, error) {
	var balance model.Balance

	query := `SELECT id, user_id, amount, updated_at FROM balances WHERE user_id = $1`
	err := r.db.Get(&balance, query, userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("balance for user ID %d not found", userID)
		}
		return nil, fmt.Errorf("error fetching balance: %w", err)
	}

	return &balance, nil
}

// UpdateAmount updates a user's balance with the new amount
func (r *BalanceRepository) UpdateAmount(userID uint64, newAmount float64) error {
	query := `UPDATE balances SET amount = $1, updated_at = NOW() WHERE user_id = $2`
	result, err := r.db.Exec(query, newAmount, userID)

	if err != nil {
		return fmt.Errorf("error updating balance: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking affected rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("no balance record found for user ID %d", userID)
	}

	return nil
}
