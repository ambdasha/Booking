package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)


//для логирования HTTP-запросов
func Logging(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()
		c.Next()

		reqID, _ := c.Get("request_id")

		logger.Info("request",
			"method", c.Request.Method,
			"path", c.FullPath(),
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
			"request_id", reqID,
		)
	}

}