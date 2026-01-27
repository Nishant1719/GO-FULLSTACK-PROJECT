package users

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

// CreateUser handles POST /users
func (h *handler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("Failed to bind request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	user, err := h.service.CreateUser(c.Request.Context(), req)
	if err != nil {
		slog.Error("Failed to create user", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GetUser handles GET /users/:id
func (h *handler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	user, err := h.service.GetUserByID(c.Request.Context(), id)
	if err != nil {
		slog.Error("Failed to get user", "error", err, "id", id)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// ListUsers handles GET /users
func (h *handler) ListUsers(c *gin.Context) {
	// Parse pagination parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Get users from service
	users, err := h.service.ListUsers(c.Request.Context(), limit, offset)
	if err != nil {
		slog.Error("Failed to list users", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch users",
		})
		return
	}

	// Get total count
	total, err := h.service.GetUsersCount(c.Request.Context())
	if err != nil {
		slog.Error("Failed to count users", "error", err)
		// Continue even if count fails
		total = 0
	}

	// Success response with pagination info
	c.JSON(http.StatusOK, gin.H{
		"data": users,
		"pagination": gin.H{
			"limit":  limit,
			"offset": offset,
			"total":  total,
		},
	})
}

// UpdateUser handles PATCH /users/:id
func (h *handler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("Failed to bind request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	user, err := h.service.UpdateUser(c.Request.Context(), id, req)
	if err != nil {
		slog.Error("Failed to update user", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser handles DELETE /users/:id
func (h *handler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	err = h.service.DeleteUser(c.Request.Context(), id)
	if err != nil {
		slog.Error("Failed to delete user", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete user",
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
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
