package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RespondWithError(ctx *gin.Context, status int, err string) {

	ctx.JSON(status, gin.H{"error": err})
}

func RespondWithSuccess(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, data)
}
func RespondWithCreated(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusCreated, data)
}
