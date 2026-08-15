package handler

import (
	"io"

	"github.com/SS-Sanjay-Kumar/Vigilis/internal/hub"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AnomalyLogsHandler struct {
	logger *zap.Logger
	hub    *hub.AnomalyDetectionHub
}

func NewAnomalyLogsHandler(l *zap.Logger, hub *hub.AnomalyDetectionHub) *AnomalyLogsHandler {
	return &AnomalyLogsHandler{
		logger: l,
		hub:    hub,
	}
}

func (alh *AnomalyLogsHandler) SendAnomalyLogs(c *gin.Context) {

	// mandatory headers for SSE
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	// subscribe client to the central hub
	clientChan := alh.hub.Subscribe()
	defer alh.hub.Unsubscribe(clientChan)

	c.Stream(func(w io.Writer) bool {

		select {
		// case: Client disconnected
		case <-c.Request.Context().Done():
			return false

		// case: receive new payload(anomaly logs) from the worker
		case msg, ok := <-clientChan:
			if !ok {
				return false
			}
			c.SSEvent("anomaly", msg)
			return true
		}
	})

}
