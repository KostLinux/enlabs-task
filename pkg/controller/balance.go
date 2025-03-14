package controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"enlabs-task/pkg/service"
)

// BalanceController handles balance-related requests
type BalanceController struct {
	balanceService service.BalanceServiceInterface
}

// NewBalanceController creates a new BalanceController instance
func NewBalanceController(balanceService service.BalanceServiceInterface) *BalanceController {
	return &BalanceController{
		balanceService: balanceService,
	}
}

// GetBalance handles GET /user/{userId}/balance requests
func (c *BalanceController) GetBalance(ctx *gin.Context) {
	userIDStr := ctx.Param("userId")

	// Parse and validate user ID
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil || userID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID. Must be a positive integer.",
		})
		return
	}

	// Get user balance
	balance, err := c.balanceService.GetUserBalance(userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve balance",
		})
		return
	}

	ctx.JSON(http.StatusOK, balance)
}
