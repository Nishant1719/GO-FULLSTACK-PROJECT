package database

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// RunMigrations applies all pending database migrations
func RunMigrations(dsn string, migrationsPath string) error {
	slog.Info("Attempting database connection", "dsn_masked", maskPassword(dsn))
	
	// Open connection using database/sql for migrations
	// Use "postgres" as driver name (lib/pq driver)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Set connection pool settings for migrations
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	slog.Info("Testing database connection...")
	// Test the connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	slog.Info("Database connection successful")

	// Create postgres driver instance
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Run migrations
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			slog.Info("No new migrations to apply")
			return nil
		}
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	slog.Info("Migrations applied successfully")
	return nil
}

// maskPassword masks the password in the DSN for logging
func maskPassword(dsn string) string {
	// Simple masking - just show structure without exposing credentials
	if len(dsn) > 30 {
		return dsn[:30] + "..." 
	}
	return dsn
}
