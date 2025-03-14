package service

import (
	"fmt"
	"math"

	"enlabs-task/pkg/model"
	"enlabs-task/pkg/repository"
)

// BalanceService handles business logic for user balances
type BalanceService struct {
	userRepo    repository.UserInterface
	balanceRepo repository.BalanceInterface
}

// NewBalanceService creates a new BalanceService instance
func NewBalanceService(
	userRepo repository.UserInterface,
	balanceRepo repository.BalanceInterface,
) *BalanceService {
	return &BalanceService{
		userRepo:    userRepo,
		balanceRepo: balanceRepo,
	}
}

// GetUserBalance retrieves a user's current balance
func (service *BalanceService) GetUserBalance(userID uint64) (*model.BalanceResponse, error) {
	// Verify the user exists
	exists, err := service.userRepo.Exists(userID)
	if err != nil {
		return nil, fmt.Errorf("error checking user existence: %w", err)
	}

	if !exists {
		return nil, fmt.Errorf("user with ID %d not found", userID)
	}

	// Get the user's balance
	balance, err := service.balanceRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving balance: %w", err)
	}

	// Format the balance to 2 decimal places
	roundedAmount := math.Round(balance.Amount*100) / 100
	formattedBalance := fmt.Sprintf("%.2f", roundedAmount)

	return &model.BalanceResponse{
		UserID:  userID,
		Balance: formattedBalance,
	}, nil
}

// UpdateBalance updates a user's balance (used internally by transaction service)
func (service *BalanceService) UpdateBalance(userID uint64, newAmount float64) error {
	if newAmount < 0 {
		return fmt.Errorf("balance cannot be negative")
	}

	return service.balanceRepo.UpdateAmount(userID, newAmount)
}
