package memstorage

import (
	"context"
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/models"
	"go.uber.org/zap"
	"sync"
)

type MemStorage struct {
	mx             sync.RWMutex
	storageCounter map[models.MetricName]int64
	storageGauge   map[models.MetricName]float64
	log            *zap.Logger
}

func NewMemStorage(log *zap.Logger) *MemStorage {
	storageCounter := make(map[models.MetricName]int64)
	storageGauge := make(map[models.MetricName]float64)

	return &MemStorage{
		storageCounter: storageCounter,
		storageGauge:   storageGauge,
		log:            log,
	}
}

func (m *MemStorage) SetMetric(ctx context.Context, metric *models.Metrics) (*models.Metrics, error) {
	if metric == nil {
		return nil, models.ErrBadBody
	}
	switch metric.MType {
	case models.MetricTypeCounter:
		m.mx.Lock()
		m.storageCounter[metric.ID] += *metric.Delta
		m.mx.Unlock()
	case models.MetricTypeGauge:
		m.mx.Lock()
		m.storageGauge[metric.ID] = *metric.Value
		m.mx.Unlock()
	}
	result, err := m.GetMetricValue(ctx, metric)
	if err != nil {
		m.log.Error(err.Error())
		return nil, err
	}
	return result, nil

}

func (m *MemStorage) GetMetricValue(ctx context.Context, metric *models.Metrics) (*models.Metrics, error) {
	m.mx.RLock()
	defer m.mx.RUnlock()
	result := &models.Metrics{
		ID:    metric.ID,
		MType: metric.MType,
	}

	switch metric.MType {
	case models.MetricTypeCounter:
		val, ok := m.storageCounter[metric.ID]
		if !ok {
			m.log.Error(models.ErrMetricNameNotFound.Error())
			return nil, models.ErrMetricNameNotFound
		}

		result.Delta = &val

	case models.MetricTypeGauge:
		val, ok := m.storageGauge[metric.ID]
		if !ok {
			m.log.Error(models.ErrMetricNameNotFound.Error())
			return nil, models.ErrMetricNameNotFound
		}

		result.Value = &val

	default:
		m.log.Error(models.ErrBadMetricType.Error())
		return nil, fmt.Errorf("%w: %s", models.ErrBadMetricType, metric.MType)
	}

	return result, nil
}

func (m *MemStorage) GetAllMetrics(ctx context.Context) ([]models.Metrics, error) {
	m.mx.RLock()
	defer m.mx.RUnlock()
	result := make([]models.Metrics, 0, len(m.storageGauge)+len(m.storageCounter))

	for k, v := range m.storageCounter {
		value := v
		m := models.Metrics{
			ID:    k,
			MType: models.MetricTypeCounter,
			Delta: &value,
		}
		result = append(result, m)
	}

	for k, v := range m.storageGauge {
		value := v
		m := models.Metrics{
			ID:    k,
			MType: models.MetricTypeGauge,
			Value: &value,
		}
		result = append(result, m)
	}

	if len(result) == 0 {
		return nil, models.ErrStorageIsEmpty
	}

	return result, nil

}

func (m *MemStorage) SetAllMetrics(ctx context.Context, metrics []models.Metrics) error {
	for i := range metrics {
		if _, err := m.SetMetric(ctx, &metrics[i]); err != nil {
			return fmt.Errorf("SetAllMetrics: %w", err)
		}
	}
	return nil
}
