package controller

import (
	"net/http"
	"time"

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
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

}
