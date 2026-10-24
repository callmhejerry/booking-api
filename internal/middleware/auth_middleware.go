package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(apiKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		x_api_key := ctx.GetHeader("X-API-Key")
		if x_api_key != apiKey {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}
	}
}
