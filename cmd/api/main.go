package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"


	_ "github.com/AbdallahZerfaoui/poweq/docs"
)

// Global logger
var logger *log.Logger

func init() {
	// Initialize logger or any other setup if needed
	logger = log.New(os.Stdout, "[API] ", log.LstdFlags)
	// Set Gin to release mode (optional)
	gin.SetMode(gin.ReleaseMode)

}

// @title Power Equation Solver API
// @version 1.0
// @description This is a REST API for solving power equations of the form x^n = K * m^x.
// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email zerfaouiabdallah@gmail.com

func main() {
	router := gin.Default()

	router.GET("/healthz", healthHandler)
	router.POST("/solve", solveHandler)

	// Swagger docs at /docs
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Set the port from environment variable or default to 8080
	port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

	// Start the server
	router.Run(":" + port)
}
