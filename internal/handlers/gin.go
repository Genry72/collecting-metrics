package handlers

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) RunServer(hostPort string) error {
	gin.SetMode(gin.DebugMode)

	g := gin.New()
	g.Use(h.ResponseLogger())
	g.Use(h.RequestLogger())

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
