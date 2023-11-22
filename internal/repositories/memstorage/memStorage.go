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

// SetMetric добавление/изменение метрики в хранилище
func (m *MemStorage) SetMetric(ctx context.Context, metrics ...*models.Metric) error {

	for i := range metrics {
		metric := metrics[i]

		if err := checkContext(ctx); err != nil {
			return fmt.Errorf("checkContext: %w", err)
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

	}

	return nil
}

// GetMetricValue получение значения метрики
func (m *MemStorage) GetMetricValue(ctx context.Context, metricType models.MetricType, metricName models.MetricName) (*models.Metric, error) {
	m.mx.RLock()
	defer m.mx.RUnlock()
	var result *models.Metric

	switch metricType {
	case models.MetricTypeCounter:
		if err := checkContext(ctx); err != nil {
			return nil, fmt.Errorf("checkContext: %w", err)
		}

		val, ok := m.storageCounter[metricName]
		if !ok {
			m.log.Error(models.ErrMetricNameNotFound.Error())
			return nil, models.ErrMetricNameNotFound
		}

		result = val

	case models.MetricTypeGauge:
		if err := checkContext(ctx); err != nil {
			return nil, fmt.Errorf("checkContext: %w", err)
		}
		val, ok := m.storageGauge[metricName]
		if !ok {
			m.log.Error(models.ErrMetricNameNotFound.Error())
			return nil, models.ErrMetricNameNotFound
		}

		result = val
	default:
		return nil, fmt.Errorf("%w: %s", models.ErrBadMetricType, metricType)
	}

	return result, nil
}

// GetAllMetrics получение всех метрик
func (m *MemStorage) GetAllMetrics(ctx context.Context) ([]*models.Metric, error) {
	m.mx.RLock()
	defer m.mx.RUnlock()
	result := make([]*models.Metric, 0, len(m.storageGauge)+len(m.storageCounter))

	for _, v := range m.storageCounter {
		if err := checkContext(ctx); err != nil {
			return nil, fmt.Errorf("GcheckContext: %w", err)
		}

		result = append(result, v)
	}

	for _, v := range m.storageGauge {
		if err := checkContext(ctx); err != nil {
			return nil, fmt.Errorf("checkContext: %w", err)
		}

		result = append(result, v)
	}

	if len(result) == 0 {
		return nil, models.ErrStorageIsEmpty
	}

	return result, nil

}

// SetAllMetrics Добавление/изменение метрик
func (m *MemStorage) SetAllMetrics(ctx context.Context, metrics []*models.Metric) error {
	for i := range metrics {
		if err := m.SetMetric(ctx, metrics[i]); err != nil {
			return fmt.Errorf("SetMetric: %w", err)
		}
	}
	return nil
}

// checkContext проверяет контекст и возвращает ошибку ErrDeadlineContext, если контекст истек.
func checkContext(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return models.ErrDeadlineContext
	default:
	}
	return nil
}
