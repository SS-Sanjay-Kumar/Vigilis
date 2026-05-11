package main

import (
	"github.com/SS-Sanjay-Kumar/Vigilis/internal/handler"
	"github.com/SS-Sanjay-Kumar/Vigilis/internal/logger"

	"github.com/gin-gonic/gin"
)

func main(){
	customLogger:= logger.GetLogger() //custom config logger 
	defer customLogger.Sync()

	healthHandler := handler.NewHealthHandler(customLogger) //dependency injection here
	logHandler := handler.NewLogHandler(customLogger) //dependency injection here

	router:= gin.Default() 
	router.GET("/health", healthHandler.CheckHealth)
	router.POST("/logs", logHandler.IngestLogs)
	router.Run("localhost:8080")
}