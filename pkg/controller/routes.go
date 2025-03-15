package controller

import (
	httpstatus "enlabs-task/pkg/http"

	"github.com/gin-gonic/gin"
)

func (ctrl *Controller) RegisterRoutes(router *gin.Engine) {
	// API routes
	api := router.Group("/user")
	{
		api.GET("/:userId/balance", ctrl.Balance.Get)
		api.POST("/:userId/transaction", ctrl.Transaction.Process)
	}

	// Health check endpoint
	router.GET("/health", func(ctx *gin.Context) {
		httpstatus.OK(ctx, "UP")
	})
}
