package worker

import (
	"fmt"

	"github.com/SS-Sanjay-Kumar/Vigilis/internal/models"
	"go.uber.org/zap"
)

type LogWorkerTools struct{
	logger *zap.Logger
	logChan chan models.LogEntry
}

func NewLogWorker(l *zap.Logger, logChan chan models.LogEntry) *LogWorkerTools{
	return &LogWorkerTools{
		logger: l,
		logChan: logChan,
	}
}

func (lw *LogWorkerTools)LogWorker() {
	lw.logger.Info("🛑 Log Worker is starting...") //using emoji to easily identify it in the console

	// Creating slices to store batches of logs
	batchSize := 5 //! batchsize is set to 5, for dev purposes
	var batch []models.LogEntry

	for log:= range lw.logChan{ //consuming logs indefinitely(loop runs indefinitely)

		batch = append(batch, log)

		if len(batch)>=batchSize { //i.e check if the batch bucket is full
			fmt.Println("🟢 Batch full; Bulk Insertion!")
			
			for i, logEntry:= range batch{
				fmt.Printf("Log no: %o Content: %+v\n",i, logEntry)
			}
			batch= nil
			fmt.Println("Batch completed!, waiting for next batch!")
		}
	}
}