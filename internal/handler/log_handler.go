package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/SS-Sanjay-Kumar/Vigilis/internal/models"
)

type LogHandler struct{
	logger *zap.Logger
}

func NewLogHandler(l *zap.Logger) *LogHandler {
	return &LogHandler{logger: l}
}

func (lh *LogHandler) IngestLogs( c *gin.Context) {

	var newLogEntry models.LogEntry // we access using package names

	if err:= c.ShouldBindJSON(&newLogEntry); err!=nil { //i.e if error
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error": err.Error(),
		})
		return
	}

	fmt.Printf("%+v\n",newLogEntry) // this prints prints the values along with its fields
	// Using just Println => {info 2026-05-11T15:09:02.618+0530 handler/health_handler.go:26 Test akjdbfhkadgf}
	// Using Printf and =V verb => {Level:info Timestamp:2026-05-11T15:09:02.618+0530 Caller:handler/health_handler.go:26 Message:Test akjdbfhkadgf}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"error": nil,
	})
}

