package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func RateLimitingMiddleware(next gin.HandlerFunc, rps int, burst int) gin.HandlerFunc {
	type Client struct {
		limiter *rate.Limiter
	}
	clients := make(map[string]*Client)
	var mu sync.Mutex

	return func(ctx *gin.Context) {
		clientIp := ctx.ClientIP()

		mu.Lock()
		client, found := clients[clientIp]
		if !found {
			client = &Client{
				limiter: rate.NewLimiter(rate.Limit(rps), burst),
			}
			clients[clientIp] = client
		}
		mu.Unlock()

		if !client.limiter.Allow() {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
			})
			return
		}
		next(ctx)
	}
}
