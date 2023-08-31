package handlers

import (
	"encoding/json"
	"github.com/Genry72/collecting-metrics/internal/models"
	"github.com/Genry72/collecting-metrics/internal/usecases/server"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	useCases *server.Server
	log      *zap.Logger
}

func NewServer(uc *server.Server, logger *zap.Logger) *Handler {
	return &Handler{
		useCases: uc,
		log:      logger,
	}
}

func (h *Handler) setMetricsText(c *gin.Context) {
	ctx := c.Request.Context()

	metricParams := &models.Metrics{}
	if err := c.ShouldBindUri(metricParams); err != nil {
		h.log.Error(err.Error())
		c.String(http.StatusBadRequest, "%v: %v", err, models.ErrFormatURL)
		return
	}

	_, status, err := h.useCases.SetMetric(ctx, metricParams)
	if err != nil {
		h.log.Error(err.Error())
		c.String(status, err.Error())
		return
	}

	c.String(status, "set metric success")

}

func (h *Handler) setMetricsJSON(c *gin.Context) {
	ctx := c.Request.Context()

	metricParams := &models.Metrics{}

	if err := c.ShouldBindJSON(metricParams); err != nil {
		h.log.Error(err.Error())
		c.String(http.StatusBadRequest, models.ErrBadBody.Error())
		return
	}

	result, status, err := h.useCases.SetMetric(ctx, metricParams)
	if err != nil {
		h.log.Error(err.Error())
		c.JSON(status, err.Error())
		return
	}

	c.JSON(status, result)
}

func (h *Handler) getMetricText(c *gin.Context) {
	ctx := c.Request.Context()

	metricParams := &models.Metrics{}

	if err := c.ShouldBindUri(metricParams); err != nil {
		h.log.Error(err.Error())
		c.String(http.StatusBadRequest, "%v", err)
		return
	}

	val, status, err := h.useCases.GetMetricValue(ctx, metricParams)
	if err != nil {
		h.log.Error(err.Error())
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

func (h *Handler) getMetricsJSON(c *gin.Context) {
	ctx := c.Request.Context()

	metricParams := &models.Metrics{}

	if err := c.ShouldBindJSON(metricParams); err != nil {
		h.log.Error(err.Error())
		c.String(http.StatusBadRequest, models.ErrBadBody.Error())
		return
	}

	val, status, err := h.useCases.GetMetricValue(ctx, metricParams)
	if err != nil {
		h.log.Error(err.Error())
		c.String(status, err.Error())
		return
	}

	valByte, err := json.Marshal(val)
	if err != nil {
		h.log.Error(err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Content-Type", "application/json")

	c.String(status, string(valByte))
}

func (h *Handler) getAllMetrics(c *gin.Context) {
	ctx := c.Request.Context()

	c.Header("Content-Type", "text/html")

	val, status, err := h.useCases.GetAllMetrics(ctx)
	if err != nil {
		h.log.Error(err.Error())
		c.String(status, err.Error())
		return
	}

	c.String(status, "%v", val)
}

func (h *Handler) pingDatabase(c *gin.Context) {

	c.Header("Content-Type", "text/html")

	if err := h.useCases.PingDataBase(); err != nil {
		h.log.Error(err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, "database connected")
}
