package worker

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type AnomalyLogsWorker struct {
	logger      *zap.Logger
	redisClient *redis.Client
}

func NewAnomalyLogsWorker(logger *zap.Logger, redisClient *redis.Client) *AnomalyLogsWorker {
	return &AnomalyLogsWorker{
		logger:      logger,
		redisClient: redisClient,
	}
}

func (aw *AnomalyLogsWorker) StartAnomalyLogsWorker() {
	aw.logger.Info("Anomaly Log Worker is starting...")

	for{

		anomaly_logs_streams, err := aw.redisClient.XRead(context.Background(), &redis.XReadArgs{
			Streams: []string{"log:anomaly", "$"}, // here "$" means, read new entries after the func call
			Block:   0,
		}).Result()

		if err != nil {
			if err == context.Canceled {
				aw.logger.Info("Anomaly worker context canceled, shutting down...")
				return
			}
			aw.logger.Error("Error reading from Redis Stream 'log:anomaly'", zap.Error(err))
			continue
		}

		fmt.Println("[temp] anomaly_logs_streams", anomaly_logs_streams) //! temp 

		//todo: to add websockets or SSE to send this value to the frontend
		// probably going to choose SSE


	}
}
