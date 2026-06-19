package middleware

import (
	"net/http"
	"strings"

	"bitedash/internal/pkg/auth"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		const bearerPrefix = "Bearer "
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}

		userID, err := auth.ParseUserIDFromToken(strings.TrimPrefix(authHeader, bearerPrefix))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set("userID", userID.String())
		c.Next()
	}
}
