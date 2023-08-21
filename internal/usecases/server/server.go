package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/models"
	"github.com/Genry72/collecting-metrics/internal/repositories"
	"go.uber.org/zap"
	"net/http"
	"strconv"

	"strings"
)

type Server struct {
	storage repositories.Repositories
	log     *zap.Logger
}

func NewServerUc(repo repositories.Repositories, log *zap.Logger) *Server {
	return &Server{
		storage: repo,
		log:     log,
	}
}

func (uc *Server) SetMetric(ctx context.Context, metric *models.Metrics) (*models.Metrics, int, error) {

	code, err := checkMetricParams(metric, true)
	if err != nil {
		uc.log.Error(err.Error())
		return nil, code, err
	}

	result, err := uc.storage.SetMetric(ctx, metric)
	if err != nil {
		uc.log.Error(err.Error())
		status := checkError(err)
		return nil, status, err
	}
	return result, http.StatusOK, nil
}

func (uc *Server) GetMetricValue(ctx context.Context, metric *models.Metrics) (*models.Metrics, int, error) {

	code, err := checkMetricParams(metric, false)
	if err != nil {
		uc.log.Error(err.Error())
		return nil, code, err
	}

	result, err := uc.storage.GetMetricValue(ctx, metric)
	if err != nil {
		uc.log.Error(err.Error())
		status := checkError(err)
		return nil, status, err
	}

	return result, http.StatusOK, nil
}

func (uc *Server) GetAllMetrics(ctx context.Context) (string, int, error) {
	mapa, err := uc.storage.GetAllMetrics(ctx)
	if err != nil {
		uc.log.Error(err.Error())
		return "", checkError(err), err
	}

	sb := strings.Builder{}

	for k, v := range mapa {
		sb.WriteString(fmt.Sprintf("%s : %v\n", k, v))
	}

	return sb.String(), http.StatusOK, nil
}

func checkMetricParams(metric *models.Metrics, checkValue bool) (int, error) {

	if checkValue && metric.ValueText == "" && metric.Value == nil && metric.Delta == nil {
		return http.StatusBadRequest, models.ErrBadMetricValue
	}

	if metric.ValueText != "" && checkValue {
		switch metric.MType {
		case models.MetricTypeGauge:
			val, err := strconv.ParseFloat(metric.ValueText, 64)
			if err != nil {
				return http.StatusBadRequest, models.ErrParseValue
			}
			metric.Value = &val
		case models.MetricTypeCounter:
			val, err := strconv.ParseInt(metric.ValueText, 10, 64)
			if err != nil {
				return http.StatusBadRequest, models.ErrParseValue
			}
			metric.Delta = &val
		}
	}

	switch metric.MType {
	case models.MetricTypeGauge:
		if checkValue {
			if metric.Value == nil {
				return http.StatusBadRequest, models.ErrBadMetricValue
			}
		}

	case models.MetricTypeCounter:
		if checkValue {
			if metric.Delta == nil {
				return http.StatusBadRequest, models.ErrBadMetricValue
			}
		}
	default:
		return http.StatusBadRequest, models.ErrBadMetricType
	}

	return http.StatusOK, nil
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
