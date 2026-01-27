package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Nishant1719/GO-FULLSTACK-PROJECT/tree/main/go-domain/internal/middleware"
	"github.com/Nishant1719/GO-FULLSTACK-PROJECT/tree/main/go-domain/internal/users"
	"github.com/gin-gonic/gin"
)

// Building HTTP Server

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

	// Health check endpoint
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Register domain routes
		users.RegisterRoutes(v1)
		// Future domains can be registered here:
		// posts.RegisterRoutes(v1)
		// products.RegisterRoutes(v1)
	}

	return r
}

// run -> gracefull shutdown
func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Second * 30,
	}
	log.Printf("Server has started : %s", app.config.addr)
	return srv.ListenAndServe()
}

// global stucture
type application struct {
	config config
	// logger
	//db driver

}

type config struct {
	addr string // server port
	db   dbConfig
	//
}

type dbConfig struct {
	dsn string // conn string
}
