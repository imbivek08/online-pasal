package database

import (
	"context"
	"fmt"
	"time"

	"github.com/imbivek08/hamropasal/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver for database/sql
)

type Database struct {
	Pool *pgxpool.Pool
}

func New(cfg *config.Config) (*Database, error) {
	// Build the DSN connection string
	dsn, err := cfg.BuildDSN(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to build DSN: %w", err)
	}

	// Configure the connection pool
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pool config: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.MaxOpenConn)
	poolConfig.MinConns = int32(cfg.MaxIdleConn)
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = time.Minute * 30

	// Create a context with timeout for connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create the connection pool
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{
		Pool: pool,
	}, nil
}

// Close closes the database connection pool
func (db *Database) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}
