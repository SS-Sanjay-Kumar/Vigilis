package handler

import (
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AnomalyLogsHandler struct {
	logger *zap.Logger
}

func NewAnomalyLogsHandler(l *zap.Logger) *AnomalyLogsHandler {
	return &AnomalyLogsHandler{
		logger: l,
	}
}

func (alh *AnomalyLogsHandler) SendAnomalyLogs(c *gin.Context) {

	// mandatory headers for SSE
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	c.Stream(func(w io.Writer) bool {

		ticker := time.NewTicker(2 * time.Second) //! temp: 2 seconds
		defer ticker.Stop()

		select {
		// Listen to client disconnection
		case <-c.Writer.CloseNotify():
			return false 

		case t := <-ticker.C:
			c.SSEvent("message", map[string]interface{}{
				"status":    "active",
				"timestamp": t.Format(time.RFC3339),
			})
			return true 
		}
	})

}
