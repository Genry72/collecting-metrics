package memstorage

import (
	"context"
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

func (m *MemStorage) SetMetricCounter(ctx context.Context, name models.MetricName, value int64) error {
	m.mx.Lock()
	defer m.mx.Unlock()

	m.storageCounter[name] += value

	return nil
}

func (m *MemStorage) SetMetricGauge(ctx context.Context, name models.MetricName, value float64) error {
	m.mx.Lock()
	defer m.mx.Unlock()

	m.storageGauge[name] = value

	return nil
}

func (m *MemStorage) GetMetricValueCounter(ctx context.Context, name models.MetricName) (int64, error) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	val, ok := m.storageCounter[name]
	if !ok {
		return 0, models.ErrMetricNameNotFound
	}

	return val, nil
}

func (m *MemStorage) GetMetricValueGauge(ctx context.Context, name models.MetricName) (float64, error) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	val, ok := m.storageGauge[name]
	if !ok {
		return 0, models.ErrMetricNameNotFound
	}

	return val, nil
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
