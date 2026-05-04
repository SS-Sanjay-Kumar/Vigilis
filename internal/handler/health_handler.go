package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

/*
In Go, we typically group handlers into a struct. 
This allows the handler to "own" its dependencies 
(like a logger or a database connection).
*/

// HealthHandler holds our dependencies
type HealthHandler struct{
	logger *zap.Logger
}

// NewHealthHandler is a constructor that "injects" the logger
func NewHealthHandler(l *zap.Logger) *HealthHandler{
	return &HealthHandler{logger: l}
}

func (h *HealthHandler) CheckHealth(c *gin.Context){
	h.logger.Info("Health Checkpoint Hit!")
	c.IndentedJSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}