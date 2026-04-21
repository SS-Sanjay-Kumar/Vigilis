package handler

import (
	"net/http"

	"github.com/SS-Sanjay-Kumar/Vigilis/internal/services"
	"github.com/gin-gonic/gin"
)

func GetAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, services.Albums)
}
