package service

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"enlabs-task/pkg/enum"
	"enlabs-task/pkg/model"
	"enlabs-task/pkg/repository"
	"enlabs-task/pkg/round"

	"gorm.io/gorm"
)

// TransactionService defines the transaction service interface for the controller
type TransactionInterface interface {
	ProcessTransaction(userID uint64, req *model.TransactionRequest, sourceType enum.SourceType) (*model.TransactionResponse, error)
}

// TransactionService handles business logic for processing transactions
type TransactionService struct {
	transactionRepo repository.TransactionInterface
	balanceRepo     repository.BalanceInterface
	userRepo        repository.UserInterface
	db              *gorm.DB
}

// NewTransactionService creates a new TransactionService instance
func NewTransactionService(
	transactionRepo repository.TransactionInterface,
	balanceRepo repository.BalanceInterface,
	userRepo repository.UserInterface,
	db *gorm.DB,
) *TransactionService {
	return &TransactionService{
		transactionRepo: transactionRepo,
		balanceRepo:     balanceRepo,
		userRepo:        userRepo,
		db:              db,
	}
}

// ProcessTransaction handles the transaction processing logic
func (service *TransactionService) ProcessTransaction(userID uint64, req *model.TransactionRequest, sourceType enum.SourceType) (*model.TransactionResponse, error) {
	var response *model.TransactionResponse

	err := service.db.Transaction(func(tx *gorm.DB) error {

		// Lock user balance row for update
		balance := &model.Balance{}
		if err := tx.Set("gorm:for_update", true).Where("user_id = ?", userID).First(balance).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("user not found")
			}

			return fmt.Errorf("error fetching balance: %w", err)
		}

		// Check for existing transaction (idempotency) with locking
		existingTx := &model.Transaction{}
		err := tx.Set("gorm:for_update", true).
			Where("transaction_id = ?", req.TransactionID).
			First(existingTx).Error

		// Handle database errors first
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("error checking transaction existence: %w", err)
		}

		// Return existing transaction if found
		if err == nil {
			response = &model.TransactionResponse{
				Success:     true,
				UserID:      userID,
				Balance:     fmt.Sprintf("%.2f", balance.Amount),
				Transaction: existingTx.TransactionID,
				ProcessedAt: existingTx.CreatedAt.Format(time.RFC3339),
			}

			return nil
		}

		// Parse and validate amount
		amount, err := strconv.ParseFloat(req.Amount, 64)
		if err != nil {
			return errors.New("invalid amount format")
		}
		amount = round.TwoDecimals(amount)

		newBalance := balance.Amount
		if req.State != "win" && balance.Amount < amount {
			return errors.New("insufficient balance")
		}

		// Apply transaction amount if a win
		if req.State == "win" {
			newBalance += amount
			return nil
		}

		// Apply transaction amount if not a win
		newBalance -= amount

		// Create transaction record
		transaction := &model.Transaction{
			TransactionID:   req.TransactionID,
			UserID:          userID,
			State:           req.State,
			Amount:          amount,
			SourceType:      string(sourceType),
			PreviousBalance: balance.Amount,
			NewBalance:      newBalance,
		}

		if err := tx.Create(transaction).Error; err != nil {
			return fmt.Errorf("error creating transaction: %w", err)
		}

		// Update balance
		if err := tx.Model(balance).Update("amount", newBalance).Error; err != nil {
			return fmt.Errorf("error updating balance: %w", err)
		}

		// Set response before completing transaction
		response = &model.TransactionResponse{
			Success:     true,
			UserID:      userID,
			Balance:     fmt.Sprintf("%.2f", newBalance),
			Transaction: req.TransactionID,
			ProcessedAt: transaction.CreatedAt.Format(time.RFC3339),
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}
