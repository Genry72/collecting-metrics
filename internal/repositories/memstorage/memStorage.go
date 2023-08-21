package memstorage

import (
	"context"
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/models"
	"sync"
)

type MemStorage struct {
	mx             sync.RWMutex
	storageCounter map[models.MetricName]int64
	storageGauge   map[models.MetricName]float64
}

func NewMemStorage() *MemStorage {
	storageCounter := make(map[models.MetricName]int64)
	storageGauge := make(map[models.MetricName]float64)

	return &MemStorage{
		storageCounter: storageCounter,
		storageGauge:   storageGauge,
	}
}

func (m *MemStorage) SetMetric(ctx context.Context, metric *models.Metrics) (*models.Metrics, error) {
	switch metric.MType {
	case string(models.MetricTypeCounter):
		m.mx.Lock()
		m.storageCounter[models.MetricName(metric.ID)] += *metric.Delta
		m.mx.Unlock()
	case string(models.MetricTypeGauge):
		m.mx.Lock()
		m.storageGauge[models.MetricName(metric.ID)] = *metric.Value
		m.mx.Unlock()
	}
	result, err := m.GetMetricValue(ctx, metric)
	if err != nil {
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
	case string(models.MetricTypeCounter):
		val, ok := m.storageCounter[models.MetricName(metric.ID)]
		if !ok {
			return nil, models.ErrMetricNameNotFound
		}

		result.Delta = &val

	case string(models.MetricTypeGauge):
		val, ok := m.storageGauge[models.MetricName(metric.ID)]
		if !ok {
			return nil, models.ErrMetricNameNotFound
		}

		result.Value = &val

	default:
		return nil, fmt.Errorf("%w: %s", models.ErrBadMetricType, metric.MType)
	}

	return result, nil
}

func (m *MemStorage) GetAllMetrics(ctx context.Context) (map[models.MetricName]interface{}, error) {
	m.mx.RLock()
	defer m.mx.RUnlock()
	result := make(map[models.MetricName]interface{}, len(m.storageGauge)+len(m.storageCounter))

	for k, v := range m.storageCounter {
		result[k] = v
	}

	for k, v := range m.storageGauge {
		result[k] = v
	}

	if len(result) == 0 {
		return nil, models.ErrStorageIsEmpty
	}

	return result, nil

}
