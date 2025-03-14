package service

import (
	"enlabs-task/pkg/model"
	"enlabs-task/pkg/repository"
)

// BalanceServiceInterface defines the contract for balance operations
type BalanceServiceInterface interface {
	GetUserBalance(userID uint64) (*model.BalanceResponse, error)
	UpdateBalance(userID uint64, newAmount float64) error
}

// TransactionServiceInterface defines the contract for transaction operations
type TransactionServiceInterface interface {
	ProcessTransaction(userID uint64, req *model.TransactionRequest, sourceType string) (*model.TransactionResponse, error)
}

// ServiceManager holds all service instances
type ServiceManager struct {
	Balance     BalanceServiceInterface
	Transaction TransactionServiceInterface
}

// NewServiceManager creates a new ServiceManager with initialized services
func NewServiceManager(repos *repository.RepositoryManager) *ServiceManager {
	balanceService := NewBalanceService(repos.User, repos.Balance)

	return &ServiceManager{
		Balance:     balanceService,
		Transaction: NewTransactionService(repos.User, repos.Balance, repos.Transaction, balanceService),
	}
}
