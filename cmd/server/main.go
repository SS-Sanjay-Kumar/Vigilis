package main

import (
	"net/http"

	"github.com/SS-Sanjay-Kumar/Vigilis/internal/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/albums", handler.GetAlbums)

	router.GET("/health", func (c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
		})
	})

	router.Run("localhost:8080")
}
