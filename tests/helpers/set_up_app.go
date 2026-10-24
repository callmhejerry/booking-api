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
	"github.com/golang-migrate/migrate/v4"
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

	database := db.New(pool)

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
