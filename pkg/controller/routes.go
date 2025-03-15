package controller

import (
	httpstatus "enlabs-task/pkg/http"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all routes with the Gin router
func (ctrl *Controller) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/user")
	{
		api.GET("/:userId/balance", ctrl.Balance.Get)
		api.POST("/:userId/transaction", ctrl.Transaction.Process)
	}

	// Add health check endpoint
	router.GET("/health", func(ctx *gin.Context) {
		httpstatus.OK(ctx, "UP")
	})

}
