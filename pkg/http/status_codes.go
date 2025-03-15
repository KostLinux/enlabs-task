package httpstatus

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	success             = http.StatusOK
	badRequest          = http.StatusBadRequest
	internalServerError = http.StatusInternalServerError
	notFound            = http.StatusNotFound
	conflict            = http.StatusConflict
	forbidden           = http.StatusForbidden
	unprocessableEntity = http.StatusUnprocessableEntity
)

func OK(ctx *gin.Context, data interface{}) {
	ctx.JSON(success, gin.H{
		"status":  success,
		"message": data,
	})
}

func BadRequest(ctx *gin.Context, message string) {
	ctx.JSON(badRequest, gin.H{
		"status":  badRequest,
		"message": message,
	})
}

func InternalServerError(ctx *gin.Context, message string) {
	ctx.JSON(internalServerError, gin.H{
		"status":  internalServerError,
		"message": message,
	})
}

func NotFound(ctx *gin.Context, message string) {
	ctx.JSON(notFound, gin.H{
		"status":  notFound,
		"message": message,
	})
}

func UnprocessableEntity(ctx *gin.Context, message string) {
	ctx.JSON(unprocessableEntity, gin.H{
		"status":  unprocessableEntity,
		"message": message,
	})
}

func Conflict(ctx *gin.Context, message string) {
	ctx.JSON(conflict, gin.H{
		"status":  conflict,
		"message": message,
	})
}

func Forbidden(ctx *gin.Context, message string) {
	ctx.JSON(forbidden, gin.H{
		"status":  forbidden,
		"message": message,
	})
}
