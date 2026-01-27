package main

import (
	"log/slog"
	"os"

	"github.com/Nishant1719/GO-FULLSTACK-PROJECT/tree/main/go-domain/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using environment variables")
	}

	// Setup structured logging
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Get database configuration from environment
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		slog.Error("DATABASE_URL environment variable is required")
		os.Exit(1)
	}
	
	slog.Info("Using database connection", "dsn", maskDSN(dsn))

	// Get server address
	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = ":8080" // Default to port 8080
	}

	// Run migrations (temporarily disabled - run manually for now)
	// TODO: Fix authentication issues with golang-migrate
	/*
	migrationsPath, err := filepath.Abs("./migrations")
	if err != nil {
		slog.Error("Failed to get migrations path", "error", err)
		os.Exit(1)
	}

	slog.Info("Running database migrations", "path", migrationsPath)
	if err := database.RunMigrations(dsn, migrationsPath); err != nil {
		slog.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}
	*/
	slog.Info("Skipping auto-migrations (run migrations manually with: docker exec -i go-domain-postgres psql -U postgres -d go_domain_db < migrations/000001_create_users_table.up.sql)")

	// Initialize database connection
	dbCfg := database.GetDefaultConfig(dsn)
	db, err := database.New(dbCfg)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer database.Close(db)

	// Create application configuration
	cfg := config{
		addr: addr,
		db: dbConfig{
			dsn:  dsn,
			pool: db,
		},
	}

	// Create and run application
	api := application{
		config: cfg,
	}

	slog.Info("Starting server", "address", cfg.addr)
	if err := api.run(api.mount()); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}

// maskDSN masks sensitive information in the DSN for logging
func maskDSN(dsn string) string {
	if len(dsn) > 40 {
		return dsn[:40] + "..."
	}
	return dsn
}
