package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"enlabs-task/pkg/model"
)

type BalanceInterface interface {
	GetByUserID(userID uint64) (*model.Balance, error)
	UpdateAmount(userID uint64, newAmount float64) error
}

// BalanceRepository handles operations on balance data
type BalanceRepository struct {
	db *gorm.DB
}

// NewBalanceRepository creates a new BalanceRepository instance
func NewBalanceRepository(db *gorm.DB) *BalanceRepository {
	return &BalanceRepository{db: db}
}

// GetByUserID retrieves a user's balance
func (repository *BalanceRepository) GetByUserID(userID uint64) (*model.Balance, error) {
	var balance model.Balance

	result := repository.db.Where("user_id = ?", userID).First(&balance)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("balance for user ID %d not found", userID)
		}
		return nil, fmt.Errorf("error fetching balance: %w", result.Error)
	}

	return &balance, nil
}

// UpdateAmount updates a user's balance with the new amount
func (repository *BalanceRepository) UpdateAmount(userID uint64, newAmount float64) error {
	result := repository.db.Model(&model.Balance{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"amount":     newAmount,
			"updated_at": gorm.Expr("NOW()"),
		})

	if result.Error != nil {
		return fmt.Errorf("error updating balance: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no balance record found for user ID %d", userID)
	}

	return nil
}
