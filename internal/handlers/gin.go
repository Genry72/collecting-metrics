package handlers

import (
	"github.com/Genry72/collecting-metrics/internal/handlers/midlware/gzip"
	"github.com/Genry72/collecting-metrics/internal/handlers/midlware/log"
	"github.com/gin-gonic/gin"
)

func (h *Handler) RunServer(hostPort string) error {
	gin.SetMode(gin.DebugMode)

	g := gin.New()
	g.Use(log.ResponseLogger(h.log))
	g.Use(log.RequestLogger(h.log))
	g.Use(gzip.Gzip(h.log))
	h.setupRoute(g)

	if err := g.Run(hostPort); err != nil {
		return err
	}

	return nil
}

func (h *Handler) setupRoute(g *gin.Engine) {
	g.GET("/", h.getAllMetrics)

	update := g.Group("update")
	update.POST("/", h.setMetricsJSON)
	update.POST("/:type/:name/:value", h.setMetricsText)

	value := g.Group("value")
	value.POST("/", h.getMetricsJSON)
	value.GET("/:type/:name", h.getMetricText)
}
