package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"enlabs-task/pkg/enum"
	"enlabs-task/pkg/model"
	"enlabs-task/pkg/service"
)

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

// ProcessTransaction handles POST /user/{userId}/transaction requests
func (c *TransactionController) Process(ctx *gin.Context) {
	// Parse and validate user ID
	userIDStr := ctx.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil || userID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get and validate Source-Type header
	sourceTypeStr := ctx.GetHeader("Source-Type")
	if sourceTypeStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing Source-Type header"})
		return
	}

	// Validate source type using the correct enum function
	sourceType, valid := enum.ParseSourceType(sourceTypeStr)
	if !valid {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Source-Type header"})
		return
	}

	// Parse request body
	var req model.TransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Process transaction
	response, err := c.transactionService.ProcessTransaction(userID, &req, string(sourceType))
	if err != nil {
		switch err.Error() {
		case "user not found":
			ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		case "invalid amount format":
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount format"})
		case "invalid transaction state":
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction state"})
		case "insufficient balance":
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Return successful response
	ctx.JSON(http.StatusOK, response)
}
