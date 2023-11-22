package handlers

import (
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/handlers/midlware/cryptor"
	"github.com/Genry72/collecting-metrics/internal/handlers/midlware/gzip"
	"github.com/Genry72/collecting-metrics/internal/handlers/midlware/log"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

// RunServer Функция использует фреймворк Gin для обработки HTTP-запросов и логирования.
// Входные параметры:
// - hostPort string: хост и порт, на котором будет запущен сервер.
// - password *string: пароль для доступа к серверу (необязательный параметр).
// Возвращаемое значение:
// - error: ошибка, возникающая при запуске сервера.
func (h *Handler) RunServer(hostPort string, password *string) error {
	gin.SetMode(gin.ReleaseMode)

	g := gin.New()

	g.Use(log.ResponseLogger(h.log))
	g.Use(log.RequestLogger(h.log))

	pprof.Register(g)

	h.setupRoute(g, password)

	if err := g.Run(hostPort); err != nil {
		return fmt.Errorf("g.Run: %w", err)
	}

	return nil
}

// setupRoute Установка хендлеров
func (h *Handler) setupRoute(g *gin.Engine, password *string) {
	g.Use(gzip.Gzip(h.log))
	if password != nil {
		g.Use(cryptor.CheckHashFromHeader(h.log, *password))
	}
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
