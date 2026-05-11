package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LogHandler struct{
	logger *zap.Logger
}

func NewLogHandler(l *zap.Logger) *LogHandler {
	return &LogHandler{logger: l}
}

type LogEntry struct{
	Level string `json:"level"`
	Timestamp string `json:"ts"`
	Caller string `json:"caller"`
	Message string `json:"msg"`
}

func (lh *LogHandler) IngestLogs( c *gin.Context) {

	var newLogEntry LogEntry

	if err:= c.ShouldBindJSON(&newLogEntry); err!=nil { //i.e if error
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": err.Error(),
		})
		return
	}

	fmt.Println(newLogEntry)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"error": nil,
	})
}

