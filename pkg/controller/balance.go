package controller

import (
	"enlabs-task/pkg/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// BalanceController handles balance-related requests
type BalanceController struct {
	balanceService service.BalanceInterface
}

// NewBalanceController creates a new BalanceController
func NewBalanceController(balanceService service.BalanceInterface) *BalanceController {
	return &BalanceController{
		balanceService: balanceService,
	}
}

// GetBalance handles GET /user/{userId}/balance requests
func (c *BalanceController) Get(ctx *gin.Context) {
	// Parse and validate user ID
	userIDStr := ctx.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil || userID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get balance
	balance, err := c.balanceService.GetBalance(userID)
	if err != nil {
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve balance"})
		}
		return
	}

	// Return balance
	ctx.JSON(http.StatusOK, balance)
}
