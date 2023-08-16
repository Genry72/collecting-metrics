package handlers

import (
	"errors"
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/models"
	"github.com/Genry72/collecting-metrics/internal/usecases/server"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
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

func (h *Handler) setMetrics(c *gin.Context) {
	ctx := c.Request.Context()

	var metric models.UpdateMetrics
	if err := c.ShouldBindUri(&metric); err != nil {
		c.String(http.StatusBadRequest, "%v: %v", err, models.ErrFormatURL)
		return
	}

	if err := h.useCases.SetMetric(ctx, &metric); err != nil {
		status := checkError(err)
		c.String(status, err.Error())
		return
	}
}

func (h *Handler) getMetricValue(c *gin.Context) {
	ctx := c.Request.Context()

	var metric models.GetMetrics

	if err := c.ShouldBindUri(&metric); err != nil {
		c.String(http.StatusBadRequest, "%v", err)
		log.Println(err)
		return
	}

	val, err := h.useCases.GetMetricValue(ctx, metric)
	if err != nil {
		status := checkError(err)
		c.String(status, err.Error())
		return
	}

	c.String(http.StatusOK, "%v", val)
}

func (h *Handler) getAllMetrics(c *gin.Context) {
	ctx := c.Request.Context()

	c.Header("Content-Type", "text/html")

	val, err := h.useCases.GetAllMetrics(ctx)

	if err != nil {
		status := checkError(err)

		c.String(status, err.Error())

		fmt.Println(err, c.Request.URL)

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
