package users

import "github.com/gin-gonic/gin"

// RegisterRoutes registers all user-related routes
func RegisterRoutes(router *gin.RouterGroup) {
	service := NewService()
	handler := NewHandler(service)

	router.GET("/users", handler.ListUsers)
	// Future routes can be added here:
	// router.GET("/users/:id", handler.GetUser)
	// router.POST("/users", handler.CreateUser)
	// router.PUT("/users/:id", handler.UpdateUser)
	// router.DELETE("/users/:id", handler.DeleteUser)
}
