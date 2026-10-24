package middleware

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		origin := ctx.GetHeader("Origin")
		if origin != "" {
			if slices.Contains(allowedOrigins, "*") {
				ctx.Header("Access-Control-Allow-Origin", "*")
			} else if slices.Contains(allowedOrigins, origin) {
				ctx.Header("Access-Control-Allow-Origin", origin)
				ctx.Header("Vary", "Origin")
			} else {
				ctx.AbortWithStatus(http.StatusForbidden)
				return
			}

			ctx.Header("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
			ctx.Header("Access-Control-Allow-Headers", "Content-Type, X-API-Key, Idempotency-Key")
			ctx.Header("Access-Control-Max-Age", "600")
		}

		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}
