package middleware

import (
	"strconv"

	"github.com/SS-Sanjay-Kumar/Vigilis/internal/metrics"
	"github.com/gin-gonic/gin"
)

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Process request

		status := strconv.Itoa(c.Writer.Status())
		metrics.HttpRequestsTotal.WithLabelValues(c.Request.Method, c.FullPath(), status).Inc()
	}
}
