package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/models"
	"github.com/Genry72/collecting-metrics/internal/repositories"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type Server struct {
	storage          repositories.Repositories
	permanentStorage repositories.PermanentStorage // Работа с файлом
	database         repositories.DatabaseStorage  // Работа с базой данных
	log              *zap.Logger
}

func NewServerUc(repo repositories.Repositories, permStor repositories.PermanentStorage, database repositories.DatabaseStorage, log *zap.Logger) *Server {
	return &Server{
		storage:          repo,
		permanentStorage: permStor,
		database:         database,
		log:              log,
	}
}

func (uc *Server) SetMetric(ctx context.Context, metrics ...*models.Metric) (int, error) {

	for i := range metrics {
		metric := metrics[i]
		code, err := checkMetricParams(metric, true)
		if err != nil {
			return code, fmt.Errorf("checkMetricParams: %w", err)
		}

		err = uc.storage.SetMetric(ctx, metric)
		if err != nil {
			status := checkError(err)
			return status, fmt.Errorf("uc.storage.SetMetric: %w", err)
		}

		// Пишем в файл все метрики из storage
		if uc.permanentStorage != nil && uc.permanentStorage.GetConfig().StoreInterval == 0 &&
			uc.permanentStorage.GetConfig().Enabled {
			if err := uc.SaveToPermanentStorage(ctx); err != nil {
				uc.log.Error(err.Error())
			}
		}
	}

	return http.StatusOK, nil
}

func (uc *Server) GetMetricValue(ctx context.Context, metric *models.Metric) (*models.Metric, int, error) {

	code, err := checkMetricParams(metric, false)
	if err != nil {
		return nil, code, fmt.Errorf("checkMetricParams: %w", err)
	}

	result, err := uc.storage.GetMetricValue(ctx, metric.MType, metric.ID)
	if err != nil {
		status := checkError(err)
		return nil, status, fmt.Errorf("uc.storage.GetMetricValue: %w", err)
	}

	return result, http.StatusOK, nil
}

func (uc *Server) GetAllMetrics(ctx context.Context) (map[models.MetricName]interface{}, int, error) {
	metrics, err := uc.storage.GetAllMetrics(ctx)
	if err != nil {
		return nil, checkError(err), fmt.Errorf("uc.storage.GetAllMetrics: %w", err)
	}

	m := make(map[models.MetricName]interface{}, len(metrics))

	for _, v := range metrics {
		switch v.MType {
		case models.MetricTypeCounter:
			m[v.ID] = v.Delta
		case models.MetricTypeGauge:
			m[v.ID] = v.Value
		default:
			return nil, http.StatusInternalServerError, fmt.Errorf("%w:%s", models.ErrMetricTypeNotFound, v.MType)
		}

	}
	return m, http.StatusOK, nil
}

func checkMetricParams(metric *models.Metric, checkValue bool) (int, error) {

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

	case errors.Is(err, models.ErrMetricTypeNotFound) || errors.Is(err, models.ErrMetricNameNotFound) ||
		errors.Is(err, sql.ErrNoRows):
		status = http.StatusNotFound

	default:
		status = http.StatusInternalServerError
	}
	return status
}
