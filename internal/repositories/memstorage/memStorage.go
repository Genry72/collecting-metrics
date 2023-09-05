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
	storageCounter map[models.MetricName]*models.Metric
	storageGauge   map[models.MetricName]*models.Metric
	log            *zap.Logger
}

func NewMemStorage(log *zap.Logger) *MemStorage {
	storageCounter := make(map[models.MetricName]*models.Metric)
	storageGauge := make(map[models.MetricName]*models.Metric)

	return &MemStorage{
		storageCounter: storageCounter,
		storageGauge:   storageGauge,
		log:            log,
	}
}

func (m *MemStorage) SetMetric(ctx context.Context, metrics ...*models.Metric) ([]*models.Metric, error) {
	result := make([]*models.Metric, 0, len(metrics))

	for i := range metrics {
		metric := metrics[i]

		if err := checkContext(ctx); err != nil {
			return nil, fmt.Errorf("SetMetric: %w", err)
		}

		switch metric.MType {
		case models.MetricTypeCounter:
			m.mx.Lock()

			_, ok := m.storageCounter[metric.ID]
			if !ok {
				m.storageCounter[metric.ID] = metric
			} else {
				*m.storageCounter[metric.ID].Delta += *metric.Delta
			}

			m.mx.Unlock()
		case models.MetricTypeGauge:
			m.mx.Lock()

			_, ok := m.storageGauge[metric.ID]
			if !ok {
				m.storageGauge[metric.ID] = metric
			} else {
				*m.storageGauge[metric.ID].Value = *metric.Value
			}

			m.mx.Unlock()
		}
		mm, err := m.GetMetricValue(ctx, metric.MType, metric.ID)
		if err != nil {
			m.log.Error(err.Error())
			return nil, err
		}
		result = append(result, mm)
	}

	return result, nil

}

func (m *MemStorage) GetMetricValue(ctx context.Context, metricType models.MetricType, metricName models.MetricName) (*models.Metric, error) {
	m.mx.RLock()
	defer m.mx.RUnlock()
	var result *models.Metric

	switch metricType {
	case models.MetricTypeCounter:
		if err := checkContext(ctx); err != nil {
			return nil, fmt.Errorf("GetMetricValue.Counter: %w", err)
		}

		val, ok := m.storageCounter[metricName]
		if !ok {
			m.log.Error(models.ErrMetricNameNotFound.Error())
			return nil, models.ErrMetricNameNotFound
		}

		result = val

	case models.MetricTypeGauge:
		if err := checkContext(ctx); err != nil {
			return nil, fmt.Errorf("GetMetricValue.Gauge: %w", err)
		}
		val, ok := m.storageGauge[metricName]
		if !ok {
			m.log.Error(models.ErrMetricNameNotFound.Error())
			return nil, models.ErrMetricNameNotFound
		}

		result = val

	default:
		m.log.Error(models.ErrBadMetricType.Error())
		return nil, fmt.Errorf("%w: %s", models.ErrBadMetricType, metricType)
	}

	return result, nil
}

func (m *MemStorage) GetAllMetrics(ctx context.Context) ([]*models.Metric, error) {
	m.mx.RLock()
	defer m.mx.RUnlock()
	result := make([]*models.Metric, 0, len(m.storageGauge)+len(m.storageCounter))

	for _, v := range m.storageCounter {
		if err := checkContext(ctx); err != nil {
			return nil, fmt.Errorf("GetAllMetrics.Counter: %w", err)
		}

		result = append(result, v)
	}

	for _, v := range m.storageGauge {
		if err := checkContext(ctx); err != nil {
			return nil, fmt.Errorf("GetAllMetrics.Gauge: %w", err)
		}

		result = append(result, v)
	}

	if len(result) == 0 {
		return nil, models.ErrStorageIsEmpty
	}

	return result, nil

}

func (m *MemStorage) SetAllMetrics(ctx context.Context, metrics []*models.Metric) error {
	for i := range metrics {
		if _, err := m.SetMetric(ctx, metrics[i]); err != nil {
			return fmt.Errorf("SetAllMetrics: %w", err)
		}
	}
	return nil
}

func checkContext(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return models.ErrDeadlineContext
	default:
	}
	return nil
}
