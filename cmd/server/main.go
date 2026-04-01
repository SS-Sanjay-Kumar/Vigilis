package main

import (
	"github.com/SS-Sanjay-Kumar/Vigilis/internal/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Health check route
	router.GET("/health", http.HealthCheck)

	// Run server on port 8000
	router.Run(":8000")
}
