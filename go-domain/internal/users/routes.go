package users

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RegisterRoutes registers all user-related routes
func RegisterRoutes(router *gin.RouterGroup, db *pgxpool.Pool) {
	// Create repository with database connection
	repo := NewPostgresRepository(db)
	
	// Create service with repository
	service := NewService(repo)
	
	// Create handler with service
	handler := NewHandler(service)

	// Register routes
	users := router.Group("/users")
	{
		users.GET("", handler.ListUsers)        // GET /api/v1/users
		users.GET("/:id", handler.GetUser)      // GET /api/v1/users/:id
		users.POST("", handler.CreateUser)      // POST /api/v1/users
		users.PATCH("/:id", handler.UpdateUser) // PATCH /api/v1/users/:id
		users.DELETE("/:id", handler.DeleteUser) // DELETE /api/v1/users/:id
	}
}
