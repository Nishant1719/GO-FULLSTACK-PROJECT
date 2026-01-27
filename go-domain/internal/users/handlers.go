package users

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) ListUsers(c *gin.Context) {
	// Call the service with request context
	users, err := h.service.ListUsers(c.Request.Context())
	if err != nil {
		// Log the error properly
		slog.Error("Failed to list users", "error", err)

		// Return error response with proper status code
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch users",
		})
		return // Stop execution!
	}

	// Success response
	c.JSON(http.StatusOK, users)
}

// Common HTTP Status Codes:
// Success
//http.StatusOK                    // 200 - Success
//http.StatusCreated               // 201 - Resource created
//http.StatusNoContent             // 204 - Success but no content

// Client Errors
//http.StatusBadRequest            // 400 - Invalid request
//http.StatusUnauthorized          // 401 - Not authenticated
//http.StatusForbidden             // 403 - Not authorized
//http.StatusNotFound              // 404 - Resource not found
//http.StatusConflict              // 409 - Resource conflict

// Server Errors
//http.StatusInternalServerError   // 500 - Server error
//http.StatusServiceUnavailable    // 503 - Service down
