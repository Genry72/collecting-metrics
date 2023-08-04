package handlers

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) RunServer(port string) error {
	gin.SetMode(gin.DebugMode)
	g := gin.New()
	h.setupRoute(g)

	if err := g.Run(":" + port); err != nil {
		return err
	}

	return nil
}

func (h *Handler) setupRoute(g *gin.Engine) {
	g.GET("/", h.getAllMetrics)

	update := g.Group("update")
	update.POST("/:type/:name/:value", h.setMetrics)

	value := g.Group("value")
	value.GET("/:type/:name", h.getMetricValue)
}
