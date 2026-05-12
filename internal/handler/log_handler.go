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
	logChannel chan models.LogEntry
}
// ! Why are we using a struct for a handler?
//* Answer: In Go, we want to avoid Global State at all costs because 
//* it makes testing impossible and leads to "Spaghetti Code."
//* By using a Struct, the handler becomes an Object that carries its own "toolbox" with it.

func NewLogHandler(l *zap.Logger,ch chan models.LogEntry) *LogHandler {
	return &LogHandler{
		logger: l,
		logChannel: ch,
	}
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
	// todo: Instead of printing the log entry, we should send the logs to the channel
	

	
	fmt.Printf("%+v\n",newLogEntry) // this prints prints the values along with its fields
	// Using just Println => {info 2026-05-11T15:09:02.618+0530 handler/health_handler.go:26 Test akjdbfhkadgf}
	// Using Printf and =V verb => {Level:info Timestamp:2026-05-11T15:09:02.618+0530 Caller:handler/health_handler.go:26 Message:Test akjdbfhkadgf}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"error": nil,
	})
}

