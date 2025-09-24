package main

import (
	// "errors"
	// "flag"
	// "fmt"
	// "os"
	"net/http"

	// "github.com/AbdallahZerfaoui/poweq/solver"
	"github.com/gin-gonic/gin"
)

func healthzHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}