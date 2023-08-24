package log

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strings"
	"time"
)

// RequestLogger Логирование входящих запросов
func RequestLogger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		c.Next()

		log.Info(
			"Request",
			zap.String("url", c.Request.RequestURI),
			zap.String("method", c.Request.Method),
			zap.Float64("latency in sec", time.Since(t).Seconds()),
			zap.String("headers Accept-Encoding", c.Request.Header.Get("Accept-Encoding")),
			zap.String("headers Content-Encoding", c.Request.Header.Get("Content-Encoding")),
		)
	}
}

// ResponseLogger Логирование ответов
func ResponseLogger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		log.Info(
			"Response",
			zap.Int("code", c.Writer.Status()),
			zap.Int("body size in bytes", c.Writer.Size()),
			zap.String("headers Content-Encoding", strings.Join(c.Writer.Header().Values("Content-Encoding"), " ,")),
		)
	}
}
