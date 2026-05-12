package main

import (
	"github.com/SS-Sanjay-Kumar/Vigilis/internal/handler"
	"github.com/SS-Sanjay-Kumar/Vigilis/internal/logger"
	"github.com/SS-Sanjay-Kumar/Vigilis/internal/models"
	"github.com/SS-Sanjay-Kumar/Vigilis/internal/worker"

	"github.com/gin-gonic/gin"
)

func main(){
	customLogger:= logger.GetLogger() //custom config logger 
	defer customLogger.Sync()
	logChan:= make(chan models.LogEntry, 5)

	healthHandler := handler.NewHealthHandler(customLogger) //dependency injection here
	logHandler := handler.NewLogHandler(customLogger, logChan)
	
	logWorker := worker.NewLogWorker(customLogger, logChan)

	go logWorker.LogWorker()

	router:= gin.Default() 
	v1:= router.Group("/v1")
	{
		v1.GET("/health", healthHandler.CheckHealth)
		v1.POST("/logs", logHandler.IngestLogs)
	}
	router.Run("localhost:8080")
}