package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/callmhejerry/booking-api/internal/config"
	"github.com/callmhejerry/booking-api/internal/db"
	"github.com/callmhejerry/booking-api/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		logger.Error("Failed to load .env file", slog.String("error", err.Error()))
		os.Exit(1)
	}

	config, err := config.LoadConfig()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	if err := db.RunMigrations("file://sql/migrations", config.DBConfig.MigrationDSN()); err != nil {
		logger.Error("Failed to apply database migrations", slog.String("error", err.Error()))
		os.Exit(1)
	}

	pool, err := db.NewDatabaseConnection(context.Background(), config.DBConfig.DSN())
	if err != nil {
		logger.Error("Failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	defer pool.Close()
	database := db.New(pool)

	server := server.NewServer(config, database, logger, pool)
	router := server.SetupRoutes()

	httpServer := http.Server{
		Handler:           router,
		Addr:              fmt.Sprintf(":%s", config.AppCongfig.Port),
		ReadTimeout:       time.Second * 10,
		WriteTimeout:      time.Second * 10,
		ReadHeaderTimeout: time.Second * 5,
	}

	go func() {
		logger.Info("Starting server", slog.String("Port", config.AppCongfig.Port))
		if err := httpServer.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {

			} else {
				logger.Error("Failed to start server", slog.String("error", err.Error()))
				os.Exit(1)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	defer close(quit)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	logger.Info("Shutting down server...")

	timeoutContext, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := httpServer.Shutdown(timeoutContext); err != nil {
		logger.Error("Failed to shut down server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
