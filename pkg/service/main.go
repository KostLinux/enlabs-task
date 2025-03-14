package service

import (
	"enlabs-task/pkg/repository"
)

// ServiceManager holds all services
type ServiceManager struct {
	Balance     *BalanceService
	Transaction *TransactionService
}

// NewServices creates a new service manager with all services
func NewServices(repos *repository.RepositoryManager) *ServiceManager {
	return &ServiceManager{
		Balance:     NewBalanceService(repos.Balance, repos.User),
		Transaction: NewTransactionService(repos.Transaction, repos.Balance, repos.User),
	}
}
