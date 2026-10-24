package helpers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http/httptest"
	"os"

	"github.com/callmhejerry/booking-api/internal/config"
	"github.com/callmhejerry/booking-api/internal/db"
	"github.com/callmhejerry/booking-api/internal/server"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type TestApp struct {
	TestServerUrl string
	Queries       *db.Queries
	Config        *config.Config
	Pool          *pgxpool.Pool
}

func NewTestApp(
	serverUrl string,
	queries *db.Queries,
	config *config.Config,
	pool *pgxpool.Pool,
) *TestApp {
	return &TestApp{
		TestServerUrl: serverUrl,
		Queries:       queries,
		Config:        config,
		Pool:          pool,
	}
}

func SetupApp() (*TestApp, func()) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	if err := godotenv.Load("../../.env.test"); err != nil {
		logger.Error("Failed to load .env.test file", slog.String("error", err.Error()))
		os.Exit(1)
	}

	config, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load test configuration", slog.String("error", err.Error()))
		os.Exit(1)
	}

	pool, err := db.NewDatabaseConnection(context.Background(), config.DBConfig.DSN())
	if err != nil {
		logger.Error("Failed to connect to test database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	database := db.New(pool)

	if err := db.RunMigrations("file://../../sql/migrations", config.DBConfig.MigrationDSN()); err != nil {
		logger.Error("Failed to apply migrations", slog.String("error", err.Error()))
		os.Exit(1)
	}

	apiServer := server.NewServer(config, database, logger, pool)

	router := apiServer.SetupRoutes()
	ts := httptest.NewServer(router)

	app := NewTestApp(ts.URL, database, config, pool)

	return app, func() {
		pool.Close()
		ts.Close()
	}
}

func (app *TestApp) CleanUpTestDB() {
	fmt.Println("Cleaning up tables")
	_, err := app.Pool.Exec(context.Background(), `
		TRUNCATE TABLE
			bookings,
			slots
		RESTART IDENTITY CASCADE;
	`)
	if err != nil {
		fmt.Println("failed to delete table data")
	}
}
