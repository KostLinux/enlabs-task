package controller

import (
	"enlabs-task/pkg/service"
)

// Controller holds all controllers
type Controller struct {
	Balance     BalanceInterface
	Transaction TransactionInterface
}

// NewController creates a new Controller with initialized controllers
func NewController(services *service.ServiceManager) *Controller {
	return &Controller{
		Balance:     NewBalanceController(services.Balance),
		Transaction: NewTransactionController(services.Transaction),
	}
}
