package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type IPLimiter struct {
	limiters typedSyncMap[string, *rate.Limiter]
	rps      rate.Limit
	burst    int
}

func NewIPLimiter(rps rate.Limit, burst int) *IPLimiter {
	return &IPLimiter{
		rps:   rps,
		burst: burst,
	}
}

func (il *IPLimiter) getLimiter(ip string) *rate.Limiter {
	if l, ok := il.limiters.Load(ip); ok {
		return l
	}

	l := rate.NewLimiter(il.rps, il.burst)
	actual, _ := il.limiters.LoadOrStore(ip, l)
	return actual
}

func RateLimit(il *IPLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !il.getLimiter(ip).Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
			})
			return
		}

		c.Next()
	}
}
