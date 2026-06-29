package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("http panic recovered",
					"method", c.Request.Method,
					"path", c.Request.URL.Path,
					"panic", rec,
					"stack", string(debug.Stack()),
					"request_id", c.GetString(RequestIDKey),
				)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
				})
			}
		}()

		c.Next()
	}
}
