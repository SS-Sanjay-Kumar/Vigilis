package main

import (
	"os"

	"github.com/SS-Sanjay-Kumar/Vigilis/internal/database"
	"github.com/SS-Sanjay-Kumar/Vigilis/internal/handler"
	"github.com/SS-Sanjay-Kumar/Vigilis/internal/logger"
	"github.com/SS-Sanjay-Kumar/Vigilis/internal/models"
	"github.com/SS-Sanjay-Kumar/Vigilis/internal/repository"
	"github.com/SS-Sanjay-Kumar/Vigilis/internal/worker"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func main() {

	if err := godotenv.Load(); err != nil {
		panic("No .env file found, using system environment variables")
	}

	customLogger := logger.GetLogger() //custom config logger
	defer customLogger.Sync()

	logChan := make(chan models.LogEntry, 5000) //! LOOK HERE: channel size

	healthHandler := handler.NewHealthHandler(customLogger) //dependency injection here
	logHandler := handler.NewLogHandler(customLogger, logChan)
	anomalyLogsHandler := handler.NewAnomalyLogsHandler(customLogger)

	// postgres
	dbConnHandler := database.NewPostgresSetup(customLogger)
	dbPoolConn, err := dbConnHandler.ConnectDB()
	if err != nil {
		customLogger.Error("🛑🛑🛑🛑🛑 Error Connection to Database 🛑🛑🛑🛑🛑", zap.Error(err))
		panic("🛑 Error Connecting to Database!!!")
	}
	//* defer closing the db connection pool(pgxpool)
	defer dbPoolConn.Close()

	// redis mq & stream
	//-----------------------------------------------------------------
	redisUrl, exists := os.LookupEnv("REDIS_URL")
	if !exists {
		customLogger.Error("🛑🛑🛑🛑🛑 Missing ENV Vars 🛑🛑🛑🛑🛑")
		panic("🛑 Missing ENV Vars: Error in Connecting to Redis!!!")

	}
	redisClient := redis.NewClient(&redis.Options{Addr: redisUrl})
	defer redisClient.Close()

	//-----------------------------------------------------------------

	logRepo := repository.NewLogRepository(dbPoolConn)
	logWorker := worker.NewLogWorker(customLogger, logChan, logRepo, redisClient)
	anomalyLogsWorker := worker.NewAnomalyLogsWorker(customLogger, redisClient)

	go logWorker.LogWorker()
	go anomalyLogsWorker.StartAnomalyLogsWorker()

	router := gin.Default()
	v1 := router.Group("/v1")
	{
		v1.GET("/health", healthHandler.CheckHealth)
		v1.POST("/logs", logHandler.IngestLogs)
		v1.GET("/events", anomalyLogsHandler.SendAnomalyLogs)
	}
	router.Run("localhost:8080")
}
