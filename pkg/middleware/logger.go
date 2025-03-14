package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger middleware for logging request details
func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Start timer
		start := time.Now()
		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery

		// Process request
		ctx.Next()

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		// Get status code and size
		statusCode := ctx.Writer.Status()
		dataLength := ctx.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		// Log request details
		fmt.Printf("[GIN] %s | %3d | %13v | %15s | %-7s %#v | %d bytes\n",
			end.Format("2006/01/02 - 15:04:05"),
			statusCode,
			latency,
			ctx.ClientIP(),
			ctx.Request.Method,
			path,
			dataLength,
		)
	}
}
