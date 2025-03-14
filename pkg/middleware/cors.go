package middleware

import (
	"enlabs-task/pkg/model"

	"github.com/gin-gonic/gin"
)

// CORS middleware for handling Cross-Origin Resource Sharing
func CORS(config *model.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", config.App.AllowOrigins)
		ctx.Header("Access-Control-Allow-Credentials", config.App.AllowCredentials)
		ctx.Header("Access-Control-Allow-Headers", config.App.AllowHeaders)
		ctx.Header("Access-Control-Allow-Methods", config.App.AllowMethods)

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}

		ctx.Next()
	}
}
