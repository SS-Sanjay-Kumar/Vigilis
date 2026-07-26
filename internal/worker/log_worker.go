package worker

import (
	"context"
	"encoding/json"
	"time"

	"github.com/SS-Sanjay-Kumar/Vigilis/internal/models"
	"github.com/SS-Sanjay-Kumar/Vigilis/internal/repository"
	"go.uber.org/zap"

	"github.com/redis/go-redis/v9"
)

type LogWorkerTools struct {
	logger  *zap.Logger
	logChan chan models.LogEntry
	repo    *repository.LogRepository
	redisMessageQueue *redis.Client
}

func NewLogWorker(l *zap.Logger, logChan chan models.LogEntry, r *repository.LogRepository, rmq *redis.Client) *LogWorkerTools {
	return &LogWorkerTools{
		logger:  l,
		logChan: logChan,
		repo:    r,
		redisMessageQueue: rmq,
	}
}

func (lw *LogWorkerTools) LogWorker() {
	lw.logger.Info("🛑 Log Worker is starting...") //using emoji to easily identify it in the console

	// Creating slices to store batches of logs
	batchSize := 100 //! LOOK HERE: batch size
	var batch []models.LogEntry

	tickerTimeInterval := 5 * time.Second
	ticker := time.NewTicker(tickerTimeInterval)

	for {
		select {
		case log, ok := <-lw.logChan:
			if !ok {
				lw.logger.Warn("LogChannel is Closed!")
				return
			}
			batch = append(batch, log)

			if len(batch) >= batchSize {
				lw.flushToDB(&batch, true)
				// todo call flushToMQ
				ticker.Reset(tickerTimeInterval) //resetting to 5 seconds
			}
		//case ends here
		case <-ticker.C:
			if len(batch) > 0 { //only consume when the channel is non-empty
				lw.flushToDB(&batch, false)
				// todo call flushToMQ
			}
		//case ends here
		}
	}
	// *batch = nil
}

func (lw *LogWorkerTools) flushToDB(batch *[]models.LogEntry, channelCase bool) { //we dont need this

	if channelCase {
		lw.logger.Info("🟢 Batch Threshold met; Bulk Insertion...") //again emoji for easy identification
	} else {
		lw.logger.Info("🟢 Batch Ticker Rang; Bulk Insertion...") //again emoji for easy identification
	}

	// stream data to postgres
	err := lw.repo.InsertLogBatch(context.Background(), *batch)
	//! set to context.background for now,
	//! need to change this for graceful shutdowns
	if err != nil {
		lw.logger.Warn("Error in Inserting Log Batch", zap.Error(err))
		return
	}

	*batch = nil//! move to line 63
	lw.logger.Info("LogWorker: Bulk Insertion Completed!")
}

//todo: another worker to drop these logs into redis message queue

// but for now, lets just use the same log worker(goroutine) to do both things

// todo for now: create a function to LPUSH the batches of logs into REDIS MQ

func (lw *LogWorkerTools) PushToRedisMQ(batch [] models.LogEntry, channelCase bool){
	if channelCase {
		lw.logger.Info("🟢 Batch Threshold met; Bulk Push to RedisMQ...") //again emoji for easy identification
	} else {
		lw.logger.Info("🟢 Batch Ticker Rang; Bulk Push to RedisMQ...") //again emoji for easy identification
	}

	//! batch size is now set to 100
	payload, err:= json.Marshal(batch)
	if err!=nil{
		lw.logger.Error("Error Parsing JSON while pushing log batches to Redis Message Queue!", zap.Error(err))
		return
	}
	
	err = lw.redisMessageQueue.LPush(context.Background(),"vigilis_log_message_queue", payload).Err()
	if err != nil {
		lw.logger.Error("Failed to push batch to Redis MQ", zap.Error(err))
		return
	}

	lw.logger.Info("LogWorker: Bulk Push to RedisMQ Completed!")

}
