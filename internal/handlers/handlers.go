package handlers

import (
	"encoding/json"
	"errors"
	"github.com/Genry72/collecting-metrics/internal/models"
	"github.com/Genry72/collecting-metrics/internal/usecases/server"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strconv"
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

	var metricParams models.SetMetricsText
	if err := c.ShouldBindUri(&metricParams); err != nil {
		c.String(http.StatusBadRequest, "%v: %v", err, models.ErrFormatURL)
		return
	}

	metric := models.Metrics{
		ID:    string(metricParams.Name),
		MType: string(metricParams.Type),
	}

	switch metricParams.Type {
	case models.MetricTypeGauge:
		val, err := strconv.ParseFloat(metricParams.Value, 64)
		if err != nil {
			c.String(http.StatusBadRequest, "%s: %s", models.ErrParseValue, err)
			return
		}

		metric.Value = &val

	case models.MetricTypeCounter:
		val, err := strconv.ParseInt(metricParams.Value, 10, 64)
		if err != nil {
			c.String(http.StatusBadRequest, "%s: %s", models.ErrParseValue, err)
			return
		}
		metric.Delta = &val

	default:
		c.String(http.StatusBadRequest, "%s: %s", models.ErrBadMetricType, metricParams.Type)
		return
	}

	if _, err := h.useCases.SetMetric(ctx, &metric); err != nil {
		status := checkError(err)
		c.String(status, err.Error())
		return
	}

	c.String(http.StatusOK, "set metric success")

}

func (h *Handler) setMetricsJson(c *gin.Context) {
	ctx := c.Request.Context()

	metricParams, err := checkMetricParams(c, true)
	if err != nil {
		h.log.Error(err.Error())
		return
	}

	result, err := h.useCases.SetMetric(ctx, metricParams)
	if err != nil {
		status := checkError(err)
		c.JSON(status, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) getMetricText(c *gin.Context) {
	ctx := c.Request.Context()

	var metric models.GetMetrics

	if err := c.ShouldBindUri(&metric); err != nil {
		c.String(http.StatusBadRequest, "%v", err)
		log.Println(err)
		return
	}

	metricParams := &models.Metrics{
		ID:    string(metric.Name),
		MType: string(metric.Type),
	}

	val, err := h.useCases.GetMetricValue(ctx, metricParams)
	if err != nil {
		status := checkError(err)
		c.String(status, err.Error())
		return
	}

	var result interface{}

	switch metric.Type {
	case models.MetricTypeCounter:
		result = *val.Delta
	case models.MetricTypeGauge:
		result = *val.Value

	}
	c.String(http.StatusOK, "%v", result)
}

func (h *Handler) getMetricsJson(c *gin.Context) {
	ctx := c.Request.Context()

	metricParams, err := checkMetricParams(c, false)
	if err != nil {
		h.log.Error(err.Error())
		return
	}

	val, err := h.useCases.GetMetricValue(ctx, metricParams)
	if err != nil {
		status := checkError(err)
		c.String(status, err.Error())
		return
	}

	valByte, err := json.Marshal(val)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	c.Header("Content-Type", "application/json")

	c.String(http.StatusOK, string(valByte))
}

func (h *Handler) getAllMetrics(c *gin.Context) {
	ctx := c.Request.Context()

	c.Header("Content-Type", "text/html")

	val, err := h.useCases.GetAllMetrics(ctx)

	if err != nil {
		status := checkError(err)

		c.String(status, err.Error())

		return
	}

	c.String(http.StatusOK, "%v", val)
}

func checkError(err error) int {
	var status int

	switch {
	case errors.Is(err, models.ErrBadMetricType) || errors.Is(err, models.ErrParseValue):
		status = http.StatusBadRequest
	case errors.Is(err, models.ErrMetricTypeNotFound) || errors.Is(err, models.ErrMetricNameNotFound):
		status = http.StatusNotFound
	default:
		status = http.StatusInternalServerError
	}
	return status
}

func checkMetricParams(c *gin.Context, checkValue bool) (*models.Metrics, error) {
	metricParams := &models.Metrics{}

	if err := c.ShouldBindJSON(metricParams); err != nil {
		c.String(http.StatusBadRequest, models.ErrBadBody.Error())
		return nil, err
	}

	switch metricParams.MType {
	case string(models.MetricTypeGauge):
		if checkValue {
			if metricParams.Value == nil {
				c.String(http.StatusBadRequest, models.ErrBadMetricValue.Error())
				return nil, models.ErrBadMetricValue
			}
		}

	case string(models.MetricTypeCounter):
		if checkValue {
			if metricParams.Delta == nil {
				c.String(http.StatusBadRequest, models.ErrBadMetricValue.Error())
				return nil, models.ErrBadMetricValue
			}
		}

	}
	return metricParams, nil
}
