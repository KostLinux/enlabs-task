package controller

import (
	"github.com/gin-gonic/gin"

	"enlabs-task/pkg/service"
)

// Controller holds all controllers
type Controller struct {
	Balance     *BalanceController
	Transaction *TransactionController
}

// NewController creates a new Controller with initialized controllers
func NewController(services *service.ServiceManager) *Controller {
	return &Controller{
		Balance:     NewBalanceController(services.Balance),
		Transaction: NewTransactionController(services.Transaction),
	}
}

// RegisterRoutes registers all routes with the Gin router
func (c *Controller) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/user")
	{
		api.GET("/:userId/balance", c.Balance.GetBalance)
		api.POST("/:userId/transaction", c.Transaction.CreateTransaction)
	}
}
