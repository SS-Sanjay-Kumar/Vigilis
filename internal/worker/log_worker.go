package worker

import (
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

	for log:= range lw.logChan {
		lw.logger.Info("Worker Consumed the Log!",
			zap.String("level", log.Level),
			zap.String("message", log.Message),
		)
	}
}