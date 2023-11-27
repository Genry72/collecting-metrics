package handlers

import (
	"bytes"
	"github.com/Genry72/collecting-metrics/internal/models"
	"github.com/Genry72/collecting-metrics/internal/usecases/server"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"html/template"
	"net/http"
)

// Handler структура для вызова методов сервера
type Handler struct {
	useCases *server.Server
	log      *zap.Logger
}

// NewServer создает экземпляр сервера
func NewServer(uc *server.Server, logger *zap.Logger) *Handler {
	return &Handler{
		useCases: uc,
		log:      logger,
	}
}

// setMetricsText установка значения метрики
func (h *Handler) setMetricsText(c *gin.Context) {
	ctx := c.Request.Context()

	metricParams := &models.Metric{}
	if err := c.ShouldBindUri(metricParams); err != nil {
		h.log.Error("setMetricsText", zap.Error(err))
		c.String(http.StatusBadRequest, "%v: %v", err, models.ErrFormatURL)
		return
	}

	status, err := h.useCases.SetMetric(ctx, metricParams)
	if err != nil {
		h.log.Error(" h.useCases.SetMetric", zap.Error(err))
		c.String(status, err.Error())
		return
	}

	c.Status(http.StatusOK)

}

// setMetricJSON отправка метрик по одной
func (h *Handler) setMetricJSON(c *gin.Context) {
	ctx := c.Request.Context()

	metricParams := &models.Metric{}

	json := jsoniter.ConfigCompatibleWithStandardLibrary

	if err := json.NewDecoder(c.Request.Body).Decode(metricParams); err != nil {
		h.log.Error("ShouldBindJSON", zap.Error(err))
		c.String(http.StatusBadRequest, models.ErrBadBody.Error())
		return
	}

	status, err := h.useCases.SetMetric(ctx, metricParams)
	if err != nil {
		h.log.Error("h.useCases.SetMetric", zap.Error(err))
		c.JSON(status, err.Error())
		return
	}

	c.Status(status)
}

// setMetricsJSON отправка метрик списком
func (h *Handler) setMetricsJSON(c *gin.Context) {
	ctx := c.Request.Context()

	metricParams := []*models.Metric{}

	json := jsoniter.ConfigCompatibleWithStandardLibrary

	if err := json.NewDecoder(c.Request.Body).Decode(&metricParams); err != nil {
		h.log.Error("ShouldBindJSON", zap.Error(err))

		if err := c.AbortWithError(http.StatusBadRequest, models.ErrBadBody); err != nil {
			h.log.Error("c.AbortWithError", zap.Error(err))
		}
		return
	}

	status, err := h.useCases.SetMetric(ctx, metricParams...)

	if err != nil {
		h.log.Error("h.useCases.SetMetric", zap.Error(err))

		if err := c.AbortWithError(status, models.ErrBadBody); err != nil {
			h.log.Error("c.AbortWithError", zap.Error(err))
		}

		return
	}

	c.Status(http.StatusOK)
}

// getMetricText получение значения метрики в текстовом формате
func (h *Handler) getMetricText(c *gin.Context) {
	ctx := c.Request.Context()

	metricParams := &models.Metric{}

	if err := c.ShouldBindUri(metricParams); err != nil {
		h.log.Error("ShouldBindUri", zap.Error(err))
		c.String(http.StatusBadRequest, "%v", err)
		return
	}

	val, status, err := h.useCases.GetMetricValue(ctx, metricParams)
	if err != nil {
		h.log.Error("h.useCases.GetMetricValue", zap.Error(err))
		c.String(status, err.Error())
		return
	}

	var result interface{}

	switch metricParams.MType {
	case models.MetricTypeCounter:
		result = *val.Delta
	case models.MetricTypeGauge:
		result = *val.Value

	}
	c.String(http.StatusOK, "%v", result)
}

// getMetricsJSON получение значения метрики
func (h *Handler) getMetricsJSON(c *gin.Context) {
	ctx := c.Request.Context()

	metricParams := &models.Metric{}

	if err := c.ShouldBindJSON(metricParams); err != nil {
		h.log.Error("c.ShouldBindJSON", zap.Error(err))
		c.String(http.StatusBadRequest, models.ErrBadBody.Error())
		return
	}

	val, status, err := h.useCases.GetMetricValue(ctx, metricParams)
	if err != nil {
		h.log.Error("h.useCases.GetMetricValue", zap.Error(err))
		c.String(status, err.Error())
		return
	}

	c.JSON(status, val)
}

// Шаблон для возврата html функции getAllMetrics
var tmpl = template.Must(template.New("metrics").Parse(`
	<!DOCTYPE html>
	<html>
	<head>
		<title>all metrics</title>
	</head>
	<body>
		<h1>all metrics</h1>
		<ul>
			{{range $key, $value := .}}
			<li>
				<strong>{{$key}}:</strong> {{$value}}
			</li>
			{{end}}
		</ul>
	</body>
	</html>
`))

// getAllMetrics Получение всех метрик в формате html
func (h *Handler) getAllMetrics(c *gin.Context) {
	ctx := c.Request.Context()

	c.Header("Content-Type", "text/html")

	val, status, err := h.useCases.GetAllMetrics(ctx)
	if err != nil {
		h.log.Error("h.useCases.GetAllMetrics", zap.Error(err))
		c.String(status, err.Error())
		return
	}

	b := bytes.Buffer{}

	err = tmpl.Execute(&b, val)
	if err != nil {
		panic(err)
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", b.Bytes())
}

// pingDatabase Проверка доступности базы данных
func (h *Handler) pingDatabase(c *gin.Context) {

	c.Header("Content-Type", "text/html")

	if err := h.useCases.PingDataBase(); err != nil {
		h.log.Error("h.useCases.PingDataBase", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, "database connected")
}
