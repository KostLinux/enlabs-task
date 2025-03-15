package controller

import (
	httpstatus "enlabs-task/pkg/http"
	"enlabs-task/pkg/model/swagger"
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

// Fix issue with swagger not being used
var _ = swagger.BalanceErrorResponse{}

// @Summary		Get user balance
// @Description	Retrieves the current balance for a specific user
// @Tags			Balance
// @Accept			json
// @Produce		json
// @Param			userId	path		int								true	"User ID"	minimum(1)
// @Success		200		{object}	model.BalanceResponse			"Successful balance retrieval"
// @Failure		400		{object}	swagger.InvalidUserIDError		"Invalid user ID provided"
// @Failure		404		{object}	swagger.UserNotFoundError		"User not found"
// @Failure		500		{object}	swagger.BalanceErrorResponse	"Internal server error"
// @Router			/user/{userId}/balance [get]
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
