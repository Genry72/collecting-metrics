package handlers

import (
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/handlers/midlware/gzip"
	"github.com/Genry72/collecting-metrics/internal/handlers/midlware/log"
	"github.com/gin-gonic/gin"
)

func (h *Handler) RunServer(hostPort string) error {
	gin.SetMode(gin.ReleaseMode)

	g := gin.New()
	g.Use(log.ResponseLogger(h.log))
	g.Use(log.RequestLogger(h.log))
	g.Use(gzip.Gzip(h.log))
	h.setupRoute(g)

	if err := g.Run(hostPort); err != nil {
		return fmt.Errorf("g.Run: %w", err)
	}

	return nil
}

func (h *Handler) setupRoute(g *gin.Engine) {
	g.GET("/", h.getAllMetrics)
	g.GET("/ping", h.pingDatabase)

	update := g.Group("update")
	update.POST("/", h.setMetricJSON)
	update.POST("/:type/:name/:value", h.setMetricsText)

	updates := g.Group("updates")
	updates.POST("/", h.setMetricsJSON)

	value := g.Group("value")
	value.POST("/", h.getMetricsJSON)
	value.GET("/:type/:name", h.getMetricText)
}
