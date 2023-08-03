package handlers

import (
	"errors"
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/models"
	"github.com/Genry72/collecting-metrics/internal/usecases"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Handler struct {
	useCases *usecases.ServerUc
}

func NewServer(uc *usecases.ServerUc) *Handler {
	return &Handler{
		useCases: uc,
	}
}

func (h *Handler) setMetrics(c *gin.Context) {
	var metric models.UpdateMetrics
	if err := c.ShouldBindUri(&metric); err != nil {
		c.String(http.StatusBadRequest, "%v: %v", err, models.ErrFormatURL)
		return
	}

	if err := h.useCases.SetMetric(&metric); err != nil {
		var status int
		if errors.Is(err, models.ErrBadMetricType) || errors.Is(err, models.ErrParseValue) {
			status = http.StatusBadRequest
		} else {
			status = http.StatusInternalServerError
		}

		c.String(status, err.Error())

		return
	}

}

func (h *Handler) getMetricValue(c *gin.Context) {
	var metric models.GetMetrics
	if err := c.ShouldBindUri(&metric); err != nil {
		c.String(http.StatusBadRequest, "%v", err)

		log.Println(err)

		return
	}

	val, err := h.useCases.GetMetricValue(metric)

	if err != nil {
		var status int
		if errors.Is(err, models.ErrMetricTypeNotFound) || errors.Is(err, models.ErrMetricNameNotFound) {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}

		c.String(status, err.Error())

		return
	}

	c.String(http.StatusOK, "%v", val)
}

func (h *Handler) getAllMetrics(c *gin.Context) {
	c.Header("Content-Type", "text/html")

	val, err := h.useCases.GetAllMetrics()

	if err != nil {
		var status int
		if errors.Is(err, models.ErrMetricTypeNotFound) || errors.Is(err, models.ErrMetricNameNotFound) {
			status = http.StatusNotFound
		} else {
			status = http.StatusInternalServerError
		}

		c.String(status, err.Error())

		fmt.Println(err, c.Request.URL)

		return
	}

	c.String(http.StatusOK, "%v", val)
}
