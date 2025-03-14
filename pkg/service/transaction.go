package service

import (
	"fmt"
	"strconv"
	"time"

	"enlabs-task/pkg/enum"
	"enlabs-task/pkg/model"
	"enlabs-task/pkg/repository"
)

// TransactionService handles business logic for transactions
type TransactionService struct {
	userRepo        repository.UserRepositoryInterface
	balanceRepo     repository.BalanceRepositoryInterface
	transactionRepo repository.TransactionRepositoryInterface
	balanceService  *BalanceService
}

// NewTransactionService creates a new TransactionService instance
func NewTransactionService(
	userRepo repository.UserRepositoryInterface,
	balanceRepo repository.BalanceRepositoryInterface,
	transactionRepo repository.TransactionRepositoryInterface,
	balanceService *BalanceService,
) *TransactionService {
	return &TransactionService{
		userRepo:        userRepo,
		balanceRepo:     balanceRepo,
		transactionRepo: transactionRepo,
		balanceService:  balanceService,
	}
}

// ProcessTransaction handles a new transaction request
func (s *TransactionService) ProcessTransaction(
	userID uint64,
	req *model.TransactionRequest,
	sourceType string,
) (*model.TransactionResponse, error) {
	// Validate user exists
	exists, err := s.userRepo.Exists(userID)
	if err != nil {
		return nil, fmt.Errorf("error checking user existence: %w", err)
	}

	if !exists {
		return nil, fmt.Errorf("user with ID %d not found", userID)
	}

	// Check for duplicate transaction (idempotency)
	existingTx, err := s.transactionRepo.FindByTransactionID(req.TransactionID)
	if err != nil {
		return nil, fmt.Errorf("error checking transaction existence: %w", err)
	}

	if existingTx != nil {
		// Transaction already processed, return the existing result
		balance, err := s.balanceRepo.GetByUserID(userID)
		if err != nil {
			return nil, fmt.Errorf("error retrieving balance: %w", err)
		}

		return &model.TransactionResponse{
			Success:     true,
			UserID:      userID,
			Balance:     fmt.Sprintf("%.2f", balance.Amount),
			Transaction: req.TransactionID,
			ProcessedAt: existingTx.CreatedAt.Format(time.RFC3339),
		}, nil
	}

	// Validate transaction state
	state, isValid := enum.ParseTransactionState(req.State)
	if !isValid {
		return nil, fmt.Errorf("invalid transaction state: %s", req.State)
	}

	// Parse amount
	amount, err := strconv.ParseFloat(req.Amount, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid amount format: %w", err)
	}

	if amount <= 0 {
		return nil, fmt.Errorf("amount must be positive")
	}

	// Get current balance
	balance, err := s.balanceRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving balance: %w", err)
	}

	// Calculate new balance based on transaction type
	previousBalance := balance.Amount
	var newBalance float64

	if state == enum.TransactionStateWin {
		newBalance = previousBalance + amount
	} else if state == enum.TransactionStateLose {
		if amount > previousBalance {
			return nil, fmt.Errorf("insufficient balance")
		}
		newBalance = previousBalance - amount
	}

	// Create transaction record
	transaction := &model.Transaction{
		TransactionID:   req.TransactionID,
		UserID:          userID,
		State:           string(state),
		Amount:          amount,
		SourceType:      sourceType,
		PreviousBalance: previousBalance,
		NewBalance:      newBalance,
	}

	// Update balance
	if err := s.balanceService.UpdateBalance(userID, newBalance); err != nil {
		return nil, fmt.Errorf("error updating balance: %w", err)
	}

	// Save transaction
	if err := s.transactionRepo.Create(transaction); err != nil {
		// Try to revert balance on error
		_ = s.balanceService.UpdateBalance(userID, previousBalance)
		return nil, fmt.Errorf("error saving transaction: %w", err)
	}

	// Return successful response
	return &model.TransactionResponse{
		Success:     true,
		UserID:      userID,
		Balance:     fmt.Sprintf("%.2f", newBalance),
		Transaction: req.TransactionID,
		ProcessedAt: transaction.CreatedAt.Format(time.RFC3339),
	}, nil
}
