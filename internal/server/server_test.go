package server

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/callmhejerry/booking-api/internal/config"
)

func TestHealthRouteIsPublic(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:password@localhost:5432/booking")
	t.Setenv("PORT", "8080")
	t.Setenv("MODE", "debug")
	t.Setenv("API_KEY", "test-key")

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	router := NewServer(cfg, nil, slog.New(slog.NewTextHandler(io.Discard, nil)), nil).SetupRoutes()
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/health", nil))

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if got, want := response.Body.String(), "{\"status\":\"ok\"}"; got != want {
		t.Errorf("body = %q, want %q", got, want)
	}
}
