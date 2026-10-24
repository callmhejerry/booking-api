package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCORSMiddlewareAllowsConfiguredPreflight(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(CORSMiddleware([]string{"http://localhost:5000"}))
	router.Use(AuthMiddleware("test-key"))
	router.POST("/bookings", func(ctx *gin.Context) {
		ctx.Status(http.StatusCreated)
	})

	request := httptest.NewRequest(http.MethodOptions, "/bookings", nil)
	request.Header.Set("Origin", "http://localhost:5000")
	request.Header.Set("Access-Control-Request-Method", "POST")
	request.Header.Set("Access-Control-Request-Headers", "content-type,x-api-key")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusNoContent)
	}
	if got := response.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5000" {
		t.Errorf("allow-origin = %q", got)
	}
}
