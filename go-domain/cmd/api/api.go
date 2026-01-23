package main

import (
	"net/http"
	"time"

	"github.com/Nishant1719/GO-FULLSTACK-PROJECT/tree/main/go-domain/internals/middleware"
	"github.com/gin-gonic/gin"
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
	r.Use(middleware.RequestID()) // assign unique id to each request
	r.Use(middleware.RealIP())    // get real client ip
	r.Use(middleware.Logger())    // log request details
	r.Use(middleware.CORS())      // handle CORS
	r.Use(middleware.Recoverer()) // recover from panics

	r.Use(middleware.Timeout(60 * time.Second)) // set timeout for requests

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
