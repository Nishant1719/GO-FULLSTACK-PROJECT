# Go Domain API

A clean architecture REST API built with Go and Gin framework.

## Project Structure

```
go-domain/
├── cmd/
│   └── api/                    # Application entry point
│       ├── main.go            # Main function, server setup
│       └── api.go             # HTTP server configuration
├── internal/                   # Private application code
│   ├── middleware/            # HTTP middlewares
│   │   └── middleware.go
│   └── users/                 # Users domain
│       ├── handlers.go        # HTTP handlers
│       ├── service.go         # Business logic
│       └── routes.go          # Route definitions
├── go.mod
├── go.sum
└── README.md
```

## Architecture

This project follows **domain-driven design** principles:

- **cmd/api**: Application entry point and server configuration
- **internal/**: Private packages (cannot be imported by external projects)
  - **users/**: Each domain is self-contained with handlers, services, and routes
  - **middleware/**: Shared middleware functions

## Running the Application

```bash
# Install dependencies
go mod download

# Run the server
go run ./cmd/api

# The server will start on http://localhost:8080
```

## API Endpoints

### Health Check
- `GET /ping` - Health check endpoint

### Users (API v1)
- `GET /api/v1/users` - List all users

## Features

- ✅ Clean architecture with domain-based organization
- ✅ Structured logging with slog
- ✅ Request ID tracking
- ✅ CORS support
- ✅ Request timeout handling
- ✅ Panic recovery
- ✅ Real IP detection

## Adding New Domains

To add a new domain (e.g., `posts`):

1. Create a new package in `internal/posts/`
2. Add `handlers.go`, `service.go`, and `routes.go`
3. Register routes in `cmd/api/api.go`:

```go
posts.RegisterRoutes(v1)
```

## Next Steps

- [ ] Add database integration
- [ ] Implement user CRUD operations
- [ ] Add authentication
- [ ] Write tests
- [ ] Add configuration management
- [ ] Add migrations
