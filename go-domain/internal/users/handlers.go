package users

import (
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
	// call the service -> List users
	users := []string{"User1", "User2"}

	// return json in an http response
	c.JSON(http.StatusOK, users)
}
