package worker

import (
	"context"
	"time"

	"github.com/SS-Sanjay-Kumar/Vigilis/internal/models"
	"github.com/SS-Sanjay-Kumar/Vigilis/internal/repository"
	"go.uber.org/zap"
)

type LogWorkerTools struct {
	logger  *zap.Logger
	logChan chan models.LogEntry
	repo    *repository.LogRepository
}

func NewLogWorker(l *zap.Logger, logChan chan models.LogEntry, r *repository.LogRepository) *LogWorkerTools {
	return &LogWorkerTools{
		logger:  l,
		logChan: logChan,
		repo:    r,
	}
}

func (lw *LogWorkerTools) LogWorker() {
	lw.logger.Info("🛑 Log Worker is starting...") //using emoji to easily identify it in the console

	// Creating slices to store batches of logs
	batchSize := 1000 //! batchsize is set to 1000, for testing purposes
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
				lw.flush(&batch, true)
				ticker.Reset(tickerTimeInterval) //resetting to 5 seconds
			}
		//case ends here
		case <-ticker.C:
			if len(batch) > 0 { //only consume when the channel is non-empty
				lw.flush(&batch, false)
			}
		//case ends here
		}
	}
}
func (lw *LogWorkerTools) flush(batch *[]models.LogEntry, channelCase bool) { //we dont need this

	if channelCase {
		lw.logger.Info("🟢 Batch Threshold met; Bulk Insertion...") //again emoji for easy identification
	} else {
		lw.logger.Info("🟢 Batch Ticker Rang; Bulk Insertion...") //again emoji for easy identification
	}

	err := lw.repo.InsertLogBatch(context.Background(), *batch)
	//! set to context.background for now,
	//! need to change this for graceful shutdowns
	if err != nil {
		lw.logger.Warn("Error in Inserting Log Batch", zap.Error(err))
		return
	}

	*batch = nil
	lw.logger.Info("LogWorker: Bulk Insertion Completed!")
}
