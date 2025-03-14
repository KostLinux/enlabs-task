package controller

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"enlabs-task/pkg/enum"
	"enlabs-task/pkg/model"
	"enlabs-task/pkg/service"
)

// TransactionController handles transaction-related requests
type TransactionController struct {
	transactionService service.TransactionServiceInterface
}

// NewTransactionController creates a new TransactionController instance
func NewTransactionController(transactionService service.TransactionServiceInterface) *TransactionController {
	return &TransactionController{
		transactionService: transactionService,
	}
}

// CreateTransaction handles POST /user/{userId}/transaction requests
func (c *TransactionController) CreateTransaction(ctx *gin.Context) {
	userIDStr := ctx.Param("userId")

	// Parse and validate user ID
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil || userID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID. Must be a positive integer.",
		})
		return
	}

	// Get and validate source type from header
	sourceTypeStr := ctx.GetHeader("Source-Type")
	if sourceTypeStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing required header: Source-Type",
		})
		return
	}

	sourceType := enum.SourceType(sourceTypeStr)
	if !sourceType.IsValid() {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Source-Type. Allowed values: game, server, payment",
		})
		return
	}

	// Parse and validate request body
	var req model.TransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Process the transaction
	response, err := c.transactionService.ProcessTransaction(userID, &req, string(sourceType))

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		if errors.Is(err, errors.New("insufficient balance")) ||
			strings.Contains(err.Error(), "insufficient balance") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Insufficient balance for this transaction",
			})
			return
		}

		if strings.Contains(err.Error(), "invalid") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to process transaction",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
