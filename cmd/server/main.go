package main

import (
	"fmt"

	"github.com/SS-Sanjay-Kumar/Vigilis/internal/database"
	"github.com/SS-Sanjay-Kumar/Vigilis/internal/handler"
	"github.com/SS-Sanjay-Kumar/Vigilis/internal/logger"
	"github.com/SS-Sanjay-Kumar/Vigilis/internal/models"
	"github.com/SS-Sanjay-Kumar/Vigilis/internal/repository"
	"github.com/SS-Sanjay-Kumar/Vigilis/internal/worker"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

func main() {

	if err := godotenv.Load(); err != nil {
		panic("No .env file found, using system environment variables")
	}

	customLogger := logger.GetLogger() //custom config logger
	defer customLogger.Sync()

	logChan := make(chan models.LogEntry, 10000) //! channel capacity is set to 10000 for testing purposes

	healthHandler := handler.NewHealthHandler(customLogger) //dependency injection here
	logHandler := handler.NewLogHandler(customLogger, logChan)

	// postgres
	dbConnHandler := database.NewPostgresSetup(customLogger)
	dbPoolConn, err := dbConnHandler.ConnectDB()
	if err != nil {
		fmt.Println(err)
		panic("🛑 Error in Connecting to Database!!!")
	}
	//* defer closing the db connection pool(pgxpool)
	defer dbPoolConn.Close()
	logRepo := repository.NewLogRepository(dbPoolConn)
	logWorker := worker.NewLogWorker(customLogger, logChan, logRepo)

	go logWorker.LogWorker()

	router := gin.Default()
	v1 := router.Group("/v1")
	{
		v1.GET("/health", healthHandler.CheckHealth)
		v1.POST("/logs", logHandler.IngestLogs)
	}
	router.Run("localhost:8080")
}
