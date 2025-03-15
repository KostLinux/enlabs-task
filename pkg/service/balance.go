package service

import (
	"fmt"

	"enlabs-task/pkg/model"
	"enlabs-task/pkg/repository"
)

// BalanceService defines the balance service interface for the controller
type BalanceInterface interface {
	GetBalance(userID uint64) (*model.BalanceResponse, error)
}

// BalanceService handles business logic for user balances
type BalanceService struct {
	balanceRepo repository.BalanceInterface
	userRepo    repository.UserInterface
}

// NewBalanceService creates a new BalanceService instance
func NewBalanceService(balanceRepo repository.BalanceInterface, userRepo repository.UserInterface) *BalanceService {
	return &BalanceService{
		balanceRepo: balanceRepo,
		userRepo:    userRepo,
	}
}

// GetBalance retrieves a user's current balance
func (service *BalanceService) GetBalance(userID uint64) (*model.BalanceResponse, error) {
	// Check if user exists
	ok, err := service.userRepo.Exists(userID)
	if err != nil {
		return nil, fmt.Errorf("error checking user existence: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	// Get balance
	balance, err := service.balanceRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Format response with balance rounded to 2 decimal places
	return &model.BalanceResponse{
		UserID:  userID,
		Balance: fmt.Sprintf("%.2f", balance.Amount),
	}, nil
}
