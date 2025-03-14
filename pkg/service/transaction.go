package service

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"enlabs-task/pkg/enum"
	"enlabs-task/pkg/model"
	"enlabs-task/pkg/repository"
)

// TransactionService defines the transaction service interface for the controller
type TransactionInterface interface {
	ProcessTransaction(userID uint64, req *model.TransactionRequest, sourceType string) (*model.TransactionResponse, error)
}

// TransactionService handles business logic for processing transactions
type TransactionService struct {
	transactionRepo repository.TransactionInterface
	balanceRepo     repository.BalanceInterface
	userRepo        repository.UserInterface
}

// NewTransactionService creates a new TransactionService instance
func NewTransactionService(
	transactionRepo repository.TransactionInterface,
	balanceRepo repository.BalanceInterface,
	userRepo repository.UserInterface,
) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		balanceRepo:     balanceRepo,
		userRepo:        userRepo,
	}
}

// ProcessTransaction handles the transaction processing logic
func (service *TransactionService) ProcessTransaction(userID uint64, req *model.TransactionRequest, sourceType string) (*model.TransactionResponse, error) {
	// Check if user exists
	exists, err := service.userRepo.Exists(userID)
	if err != nil {
		return nil, fmt.Errorf("error checking user existence: %w", err)
	}
	if !exists {
		return nil, errors.New("user not found")
	}

	// Check if transaction already exists (idempotence)
	existingTx, err := service.transactionRepo.FindByTransactionID(req.TransactionID)
	if err != nil {
		return nil, fmt.Errorf("error checking transaction existence: %w", err)
	}

	// If transaction was already processed, return the existing result
	if existingTx != nil {
		// Get current balance to return in response
		balance, err := service.balanceRepo.GetByUserID(userID)
		if err != nil {
			return nil, fmt.Errorf("error fetching balance: %w", err)
		}

		return &model.TransactionResponse{
			Success:     true,
			UserID:      userID,
			Balance:     fmt.Sprintf("%.2f", balance.Amount),
			Transaction: req.TransactionID,
			ProcessedAt: existingTx.CreatedAt.Format(time.RFC3339),
		}, nil
	}

	// Parse transaction amount
	amount, err := strconv.ParseFloat(req.Amount, 64)
	if err != nil {
		return nil, errors.New("invalid amount format")
	}

	// Round to 2 decimal places
	amount = float64(int64(amount*100)) / 100

	// Validate state
	state := enum.TransactionState(req.State)
	if !state.IsTransactionValid() {
		return nil, errors.New("invalid transaction state")
	}

	// Get current balance
	currentBalance, err := service.balanceRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching balance: %w", err)
	}

	// Calculate new balance based on transaction state
	var newBalance float64
	if state == enum.TransactionStateWin {
		newBalance = currentBalance.Amount + amount
	} else {
		// Check if sufficient balance for 'lose' transaction
		if currentBalance.Amount < amount {
			return nil, errors.New("insufficient balance")
		}
		newBalance = currentBalance.Amount - amount
	}

	// Round the new balance to 2 decimal places
	newBalance = float64(int64(newBalance*100)) / 100

	// Create transaction record
	tx := &model.Transaction{
		TransactionID:   req.TransactionID,
		UserID:          userID,
		State:           string(state),
		Amount:          amount,
		SourceType:      sourceType,
		PreviousBalance: currentBalance.Amount,
		NewBalance:      newBalance,
	}

	// Save transaction
	if err := service.transactionRepo.Create(tx); err != nil {
		return nil, fmt.Errorf("error creating transaction: %w", err)
	}

	// Update user balance
	if err := service.balanceRepo.UpdateAmount(userID, newBalance); err != nil {
		return nil, fmt.Errorf("error updating balance: %w", err)
	}

	// Return response
	return &model.TransactionResponse{
		Success:     true,
		UserID:      userID,
		Balance:     fmt.Sprintf("%.2f", newBalance),
		Transaction: req.TransactionID,
		ProcessedAt: tx.CreatedAt.Format(time.RFC3339),
	}, nil
}
