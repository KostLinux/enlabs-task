package controller

import (
	httpstatus "enlabs-task/pkg/http"
	"enlabs-task/pkg/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BalanceInterface interface {
	Get(ctx *gin.Context)
}

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
func (ctrl *BalanceController) Get(ctx *gin.Context) {
	// Parse and validate user ID
	userIDStr := ctx.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil || userID == 0 {
		httpstatus.BadRequest(ctx, "Invalid user ID")
		return
	}

	// Get balance
	balance, err := ctrl.balanceService.GetBalance(userID)
	if err != nil {
		if err.Error() == "user not found" {
			httpstatus.NotFound(ctx, "User not found")
			return
		}

		httpstatus.InternalServerError(ctx, "Failed to retrieve balance")
		return
	}

	// Return balance
	httpstatus.OK(ctx, balance)
}
