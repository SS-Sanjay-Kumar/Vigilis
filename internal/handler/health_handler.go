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
type HealthHandler struct {
	logger *zap.Logger
}

// NewHealthHandler is a constructor that "injects" the logger
func NewHealthHandler(l *zap.Logger) *HealthHandler {
	return &HealthHandler{logger: l}
}

func (h *HealthHandler) CheckHealth(c *gin.Context) {
	h.logger.Info("Health Checkpoint Hit!(temp)") // we have to push this into the channel
	// eg: {
	// "level":"info",
	// "ts":"2026-05-11T15:09:02.618+0530",
	// "caller":"handler/health_handler.go:26",
	// "msg":"Health Checkpoint Hit!"
	// }
	c.IndentedJSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// !todo: Implement concurrent buffer(Producer Consumer Pattern)
// Instead of logging a string(refer l.no 28, we should log a structured type)
