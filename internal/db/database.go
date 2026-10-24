package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
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

func WithTx(ctx context.Context, dbConn *pgxpool.Pool, txnFunc func(qtx *Queries) error) error {
	txn, err := dbConn.Begin(ctx)

	if err != nil {
		return err
	}

	defer txn.Rollback(ctx)

	qtx := New(dbConn)

	err = txnFunc(qtx)
	if err != nil {
		return err
	}
	txn.Commit(ctx)
	return nil
}
