package main

import (
	// "errors"
	// "flag"
	// "fmt"
	// "os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	// "github.com/yourusername/poweq/solver"
)

// @title Power Equation Solver API
// @version 1.0
// @description This is a REST API for solving power equations of the form x^n = K * m^x.
// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email zerfaouiabdallah@gmail.com

func main() {
	router := gin.Default()

	router.GET("/healthz", healthzHandler)

	// Swagger docs at /docs
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	// Start the server
	router.Run("localhost:8080")
}