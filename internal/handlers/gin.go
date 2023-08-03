package handlers

import (
	"github.com/Genry72/collecting-metrics/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
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

	g.NoRoute(func(c *gin.Context) {
		if c.Request.Method == http.MethodPost {
			return
		}
		c.String(http.StatusMethodNotAllowed, models.ErrOnlyPost.Error())
	})

}
