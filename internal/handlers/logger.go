package handlers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

// RequestLogger Логирование входящих запросов
func (h *Handler) RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		c.Next()

		h.log.Info(
			"Request",
			zap.String("url", c.Request.RequestURI),
			zap.String("method", c.Request.Method),
			zap.Float64("latency in sec", time.Since(t).Seconds()),
			//zap.Int("code", c.Request.Response.StatusCode),
		)
	}
}

// ResponseLogger Логирование ответов
func (h *Handler) ResponseLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		h.log.Info(
			"Response",
			zap.Int("code", c.Writer.Status()),
			zap.Int("body size in bytes", c.Writer.Size()),
		)
	}
}
