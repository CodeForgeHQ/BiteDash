package middleware

import (
	"time"

	"bitedash/internal/metrics"

	"github.com/gin-gonic/gin"
)

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startedAt := time.Now()

		c.Next()

		metrics.ObserveHTTPRequest(c, time.Since(startedAt))
	}
}
