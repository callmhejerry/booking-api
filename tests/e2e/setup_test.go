package e2e

import (
	"context"
	"fmt"
	"log/slog"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/callmhejerry/booking-api/internal/config"
	"github.com/callmhejerry/booking-api/internal/db"
	"github.com/callmhejerry/booking-api/internal/server"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var TestServerURL string
var Queries db.Queries

func TestMain(m *testing.M) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	err := godotenv.Load("../../.env.test")
	if err != nil {
		logger.Error("Failed to load .env.test file", slog.String("error", err.Error()))
		os.Exit(1)
	}

	config := config.LoadConfig()

	pool, err := db.NewDatabaseConnection(context.Background(), config.DBConfig.DSN())
	if err != nil {
		logger.Error("Failed to connect to test database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	defer pool.Close()
	database := db.New(pool)
	Queries = *database

	migrator, err := migrate.New(
		"file://../../sql/migrations", fmt.Sprintf("pgx5://%s:%s@%s:%s/%s?sslmode=%s", config.DBConfig.Username, config.DBConfig.Password, config.DBConfig.Host, config.DBConfig.Port, config.DBConfig.Name, config.DBConfig.SSLMode),
	)
	if err != nil {
		logger.Error("Failed to create migrator",
			"error", err.Error(),
		)
		os.Exit(1)
	}

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Error("Failed to apply migration")
		os.Exit(1)
	}

	server := server.NewServer(config, database, logger, pool)

	router := server.SetupRoutes()

	ts := httptest.NewServer(router)
	defer ts.Close()

	TestServerURL = ts.URL

	code := m.Run()
	cleanUpTestDB(pool)
	os.Exit(code)
}

func cleanUpTestDB(pool *pgxpool.Pool) {
	_, err := pool.Exec(context.Background(), `
		TRUNCATE TABLE
			bookings,
			slots
		RESTART IDENTITY CASCADE;
	`)
	if err != nil {
		fmt.Println("failed to delete table data")
	}
}
