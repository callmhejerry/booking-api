package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func RateLimitingMiddleware(next gin.HandlerFunc, rps int, burst int) gin.HandlerFunc {
	type Client struct {
		limiter *rate.Limiter
	}
	clients := make(map[string]*Client)

	return func(ctx *gin.Context) {
		clientIp := ctx.ClientIP()
		client, found := clients[clientIp]
		if !found {
			client = &Client{
				limiter: rate.NewLimiter(rate.Limit(rps), burst),
			}
			clients[clientIp] = client
		}
		if !client.limiter.Allow() {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
			})
		}
		next(ctx)
	}
}
