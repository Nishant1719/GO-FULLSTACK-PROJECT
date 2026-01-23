package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-domain/internals/middleware"
)

// Building HTTP Server
// global stucture
type application struct {
	config config
	// logger
	//db driver

}

// mount
func (app *application) mount() http.Handler {
	r := gin.Default()

	// Apply middlewares
	r.Use(middleware.RequestID())
	r.Use(middleware.RealIP())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())
	r.Use(middleware.Timeout(60 * time.Second))

	// Define a simple GET endpoint
	r.GET("/ping", func(c *gin.Context) {
		// Return JSON response
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	return r
}

// run -> gracefull shutdown
type config struct {
	addr string // server port
	db   dbConfig
	//
}

type dbConfig struct {
	dsn string // conn string
}
