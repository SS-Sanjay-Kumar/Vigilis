package worker

import (
	"context"
	"encoding/json"

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

	for {

		anomaly_logs_streams, err := aw.redisClient.XRead(context.Background(), &redis.XReadArgs{
			Streams: []string{"log:anomaly", "$"}, // here "$" means, read new entries after the func call
			Block:   0,
			Count:   10, //! temp
		}).Result()

		if err != nil {
			if err == context.Canceled {
				aw.logger.Info("Anomaly worker context canceled, shutting down...")
				return
			}
			aw.logger.Error("Error reading from Redis Stream 'log:anomaly'", zap.Error(err))
			continue
		}

		// fmt.Println("[temp] anomaly_logs_streams", anomaly_logs_streams)
		//output:
		// [temp] anomaly_logs_streams [{log:anomaly [{1786777006079-0 map[caller:dfs.DataNode$DataXceiver level:info message:sample maja bek veedu mse:0.021571 threshold:0.001085 timestamp:2008-11-11T10:26:43+05:30] 0 0}]}]
		//* visualising it
		// [temp] anomaly_logs_streams
		// [
		// 		(stream (i.e anomaly_logs_streams)){
		//			log:anomaly [
		// 					(stream.Messages){
		// 						1786774751105-0 map[caller:dfs.DataNode$DataXceiver level:info message:sample maja bek  mse:0.022578 threshold:0.001085 timestamp:2008-11-11T10:26:43+05:30] 0 0
		// note: message.Values = map(string: any)
		//					}
		//			]
		// 		}
		// ]

		for _, stream := range anomaly_logs_streams {
			// fmt.Println("‼️‼️‼️stream ", stream)
			for _, message := range stream.Messages {
				// fmt.Println("‼️‼️‼️message ", message)
				// fmt.Println("‼️‼️‼️message.Values ", message.Values)

				// message.Values is map[string]interface{}
				jsonBytes, err := json.Marshal(message.Values)
				if err != nil {
					aw.logger.Error("Failed to marshal stream payload to JSON", zap.Error(err))
					continue
				}

				jsonString := string(jsonBytes)
				aw.logger.Info("JSON -> " + jsonString)
				//* checking:
				// fmt.Println("Extracted JSON:", jsonString)
				// fmt.Println("logsBatch", logsBatch)
				// output:
				// Extracted JSON: {"caller":"dfs.DataNode$DataXceiver","level":"info","message":"sample checks logsBatch","mse":"0.022590","threshold":"0.001085","timestamp":"2026-11-11T10:26:43+05:30"}
				// logsBatch [{"caller":"dfs.DataNode$DataXceiver","level":"info","message":"sample checks logsBatch","mse":"0.022590","threshold":"0.001085","timestamp":"2026-11-11T10:26:43+05:30"}]

			}
		}

		// [temp] anomaly_logs_streams [{log:anomaly [{1786777006079-0 map[caller:dfs.DataNode$DataXceiver level:info message:sample maja bek veedu mse:0.021571 threshold:0.001085 timestamp:2008-11-11T10:26:43+05:30] 0 0}]}]
		// ‼️‼️‼️stream  {log:anomaly [{1786777006079-0 map[caller:dfs.DataNode$DataXceiver level:info message:sample maja bek veedu mse:0.021571 threshold:0.001085 timestamp:2008-11-11T10:26:43+05:30] 0 0}]}
		// ‼️‼️‼️message  {1786777006079-0 map[caller:dfs.DataNode$DataXceiver level:info message:sample maja bek veedu mse:0.021571 threshold:0.001085 timestamp:2008-11-11T10:26:43+05:30] 0 0}
		// ‼️‼️‼️message.Values  map[caller:dfs.DataNode$DataXceiver level:info message:sample maja bek veedu mse:0.021571 threshold:0.001085 timestamp:2008-11-11T10:26:43+05:30]
		// Extracted JSON: {"caller":"dfs.DataNode$DataXceiver","level":"info","message":"sample maja bek veedu","mse":"0.021571","threshold":"0.001085","timesta

		// sse: send anomaly logs to the frontend
	}

}
