package worker

import (
	"fmt"
	"time"

	"github.com/SS-Sanjay-Kumar/Vigilis/internal/models"
	"go.uber.org/zap"
)

type LogWorkerTools struct {
	logger  *zap.Logger
	logChan chan models.LogEntry
}

func NewLogWorker(l *zap.Logger, logChan chan models.LogEntry) *LogWorkerTools {
	return &LogWorkerTools{
		logger:  l,
		logChan: logChan,
	}
}

func (lw *LogWorkerTools) LogWorker() {
	lw.logger.Info("🛑 Log Worker is starting...") //using emoji to easily identify it in the console

	// Creating slices to store batches of logs
	batchSize := 5 //! batchsize is set to 5, for dev purposes
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
				lw.logger.Info("🟢 Batch Threshold met; Bulk Insertion...") //again emoji for easy identification

				for i, entry := range batch {
					msg := fmt.Sprintf("Log no: %o; Entry: %+v\n", i, entry)
					lw.logger.Info(msg)
				}
				batch = nil
				ticker.Reset(tickerTimeInterval) //resetting to 5 seconds
				lw.logger.Info("LogWorker: Bulk Insertion Completed!")
			}
		//case ends here
		case <-ticker.C:
			if len(batch) > 0 { //only consume when the channel is non-empty
				lw.logger.Info("🟢 Batch Ticker Rang; Bulk Insertion...") //again emoji for easy identification

				for i, entry := range batch {
					msg := fmt.Sprintf("Log no: %o; Entry: %+v\n", i, entry)
					lw.logger.Info(msg)
				}
				batch = nil
				lw.logger.Info("LogWorker: Bulk Insertion Completed!")
			}
			//case ends here
		}
	}
}
