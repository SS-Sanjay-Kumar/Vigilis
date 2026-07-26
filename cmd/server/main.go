package main

import (
	"fmt"
	"os"

	"github.com/SS-Sanjay-Kumar/Vigilis/internal/database"
	"github.com/SS-Sanjay-Kumar/Vigilis/internal/handler"
	"github.com/SS-Sanjay-Kumar/Vigilis/internal/logger"
	"github.com/SS-Sanjay-Kumar/Vigilis/internal/models"
	"github.com/SS-Sanjay-Kumar/Vigilis/internal/repository"
	"github.com/SS-Sanjay-Kumar/Vigilis/internal/worker"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

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

	// postgres
	dbConnHandler := database.NewPostgresSetup(customLogger)
	dbPoolConn, err := dbConnHandler.ConnectDB()
	if err != nil {
		fmt.Println(err)
		panic("🛑 Error in Connecting to Database!!!")
	}
	//* defer closing the db connection pool(pgxpool)
	defer dbPoolConn.Close()

	// redis mq
	//-----------------------------------------------------------------
	redisUrl, exists := os.LookupEnv("REDIS_URL")
	if !exists{
		fmt.Println("🛑🛑🛑🛑🛑 Missing ENV Vars 🛑🛑🛑🛑🛑")
		panic("🛑 Error in Connecting to Redis!!!")

	}
	redisMQClient := redis.NewClient(&redis.Options{Addr: redisUrl})
	fmt.Println(redisMQClient)

	//-----------------------------------------------------------------

	logRepo := repository.NewLogRepository(dbPoolConn)
	logWorker := worker.NewLogWorker(customLogger, logChan, logRepo, redisMQClient)

	go logWorker.LogWorker()

	router := gin.Default()
	v1 := router.Group("/v1")
	{
		v1.GET("/health", healthHandler.CheckHealth)
		v1.POST("/logs", logHandler.IngestLogs)
	}
	router.Run("localhost:8080")
}
