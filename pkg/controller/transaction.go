package controller

import (
	"strconv"

	"enlabs-task/pkg/model/swagger"

	"github.com/gin-gonic/gin"

	"enlabs-task/pkg/enum"
	httpstatus "enlabs-task/pkg/http"
	"enlabs-task/pkg/model"
	"enlabs-task/pkg/service"
)

// Fix issue with swagger not being used
var _ = swagger.InvalidUserIDError{}

type TransactionInterface interface {
	Process(ctx *gin.Context)
}

// TransactionController handles transaction-related requests
type TransactionController struct {
	transactionService service.TransactionInterface
}

// NewTransactionController creates a new TransactionController
func NewTransactionController(transactionService service.TransactionInterface) *TransactionController {
	return &TransactionController{
		transactionService: transactionService,
	}
}

// @Summary		Create transaction
// @Description	Process a win/lose transaction for a specific user
// @Tags			Transaction
// @Accept			json
// @Produce		json
// @Param			userId		path		int							true	"User ID"				minimum(1)
// @Param			Source-Type	header		string						true	"Transaction source"	Enums(game,server,payment)
// @Param			transaction	body		model.TransactionRequest	true	"Transaction details"
// @Success		200			{object}	model.TransactionResponse
// @Failure		400			{object}	swagger.InvalidUserIDError				"Invalid user ID"
// @Failure		400			{object}	swagger.InvalidSourceTypeError			"Invalid source type"
// @Failure		400			{object}	swagger.InvalidAmountError				"Invalid amount"
// @Failure		400			{object}	swagger.InvalidTransactionStateError	"Invalid transaction state"
// @Failure		404			{object}	swagger.UserNotFoundError				"User not found"
// @Failure		404			{object}	swagger.MissingSourceTypeError			"Missing source type"
// @Failure		422			{object}	swagger.InsufficientBalanceError		"Insufficient balance"
// @Router			/user/{userId}/transaction [post]
func (ctrl *TransactionController) Process(ctx *gin.Context) {
	// Parse and validate user ID
	userIDStr := ctx.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil || userID == 0 {
		httpstatus.BadRequest(ctx, "Invalid user ID")
		return
	}

	// Get and validate Source-Type header
	sourceTypeStr := ctx.GetHeader("Source-Type")
	if sourceTypeStr == "" {
		httpstatus.NotFound(ctx, "Source-Type header not found")
		return
	}

	// Validate source type using the correct enum function
	sourceType, valid := enum.ParseSourceType(sourceTypeStr)
	if !valid {
		httpstatus.BadRequest(ctx, "Invalid Source-Type header")
		return
	}

	// Parse request body
	var req model.TransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		httpstatus.BadRequest(ctx, "Invalid request body. Please check docs and try again.")
		return
	}

	type errorResponse struct {
		statusHandler func(*gin.Context, string)
		message       string
	}

	// Map errors to their appropriate HTTP responses
	var transactionErrors = map[string]errorResponse{
		"user not found":            {httpstatus.NotFound, "User not found"},
		"invalid amount format":     {httpstatus.BadRequest, "Invalid amount format"},
		"invalid transaction state": {httpstatus.BadRequest, "Invalid transaction state"},
		"insufficient balance":      {httpstatus.UnprocessableEntity, "Insufficient balance"},
	}

	// Process transaction
	response, err := ctrl.transactionService.ProcessTransaction(userID, &req, sourceType)
	if err != nil {
		if errorResponse, ok := transactionErrors[err.Error()]; ok {
			errorResponse.statusHandler(ctx, errorResponse.message)
			return
		}

		httpstatus.InternalServerError(ctx, "Failed to process transaction. Please try again.")
		return
	}

	// Return successful response
	httpstatus.OK(ctx, response)
}
