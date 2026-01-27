package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Config holds database configuration
type Config struct {
	DSN             string
	MaxOpenConns    int32
	MaxIdleConns    int32
	MaxIdleTime     time.Duration
	MaxConnLifetime time.Duration
}

// New creates a new database connection pool
func New(cfg Config) (*pgxpool.Pool, error) {
	// Parse connection string
	config, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database DSN: %w", err)
	}

	// Set connection pool configuration
	config.MaxConns = cfg.MaxOpenConns
	config.MinConns = cfg.MaxIdleConns
	config.MaxConnIdleTime = cfg.MaxIdleTime
	config.MaxConnLifetime = cfg.MaxConnLifetime

	// Create connection pool
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	slog.Info("Database connection established successfully")
	return pool, nil
}

// Close gracefully closes the database connection pool
func Close(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
		slog.Info("Database connection closed")
	}
}

// GetDefaultConfig returns a database configuration with sensible defaults
func GetDefaultConfig(dsn string) Config {
	return Config{
		DSN:             dsn,
		MaxOpenConns:    25,              // Maximum number of open connections
		MaxIdleConns:    5,               // Maximum number of idle connections
		MaxIdleTime:     15 * time.Minute, // Maximum idle time
		MaxConnLifetime: time.Hour,        // Maximum connection lifetime
	}
}
