package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func NewDatabaseConnection(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	dbConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return nil, fmt.Errorf("Create database pool connection: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("Failed to ping database %w", err)
	}
	return pool, nil
}

func RunMigrations(sourceURL, databaseURL string) (err error) {
	migrator, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return fmt.Errorf("create database migrator: %w", err)
	}
	defer func() {
		sourceErr, databaseErr := migrator.Close()
		err = errors.Join(err, sourceErr, databaseErr)
	}()

	if migrationErr := migrator.Up(); migrationErr != nil && !errors.Is(migrationErr, migrate.ErrNoChange) {
		return fmt.Errorf("apply database migrations: %w", migrationErr)
	}
	return nil
}

func WithTx(ctx context.Context, dbConn *pgxpool.Pool, txnFunc func(qtx *Queries) error) error {
	txn, err := dbConn.Begin(ctx)

	if err != nil {
		return err
	}

	defer txn.Rollback(ctx)

	qtx := New(dbConn).WithTx(txn)

	err = txnFunc(qtx)
	if err != nil {
		return err
	}
	return txn.Commit(ctx)
}
