package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Structs for request and response payloads
type SolveRequest struct {
	N         float64 `json:"n" binding:"required" example:"2"`
	M         float64 `json:"m" binding:"required" example:"2.718281828"`
	K         float64 `json:"k" binding:"required" example:"1"`
	A         float64 `json:"a" binding:"required" example:"0.1"`
	B         float64 `json:"b" binding:"required" example:"10"`
	Tolerance float64 `json:"tolerance" example:"0.000001"`
	MaxIter   int     `json:"max_iter" example:"100"`
	Algorithm string  `json:"algorithm" example:"newton"`
}

type SolveResponse struct {
	Solutions []APISolution `json:"solutions"`
}

type APISolution struct {
	X     float64 `json:"x" example:"2.0"`
	Steps int     `json:"steps" example:"5"`
	Error error   `json:"error,omitempty"`
}

// healthHandler handles the health check endpoint.
// @Summary Health check
// @Description Returns the API status
// @Tags health
// @Produce  json
// @Success 200 {object} map[string]string
// @Router /healthz [get]
func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Solve godoc
// @Summary Solve a power equation
// @Description Solves x^n = K * m^x using the specified algorithm
// @Tags solver
// @Accept  json
// @Produce  json
// @Param   request body SolveRequest true "Solve Request"
// @Success 200 {object} SolveResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /solve [post]
func solveHandler(c *gin.Context) {
	var req SolveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the solver function (to be implemented)
	result, err := req.Solve4API()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}
