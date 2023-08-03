package repositories

import (
	"github.com/Genry72/collecting-metrics/internal/models"
	"sync"
)

type MemStorage struct {
	sync.RWMutex
	storage map[models.MetricType]map[models.MetricName]interface{}
}

func NewMemStorage() *MemStorage {
	storage := make(map[models.MetricType]map[models.MetricName]interface{})

	return &MemStorage{
		storage: storage,
	}
}

func (m *MemStorage) SetMetric(metric Metric) error {
	m.Lock()
	defer m.Unlock()

	if m.storage[metric.GetType()] == nil {
		m.storage[metric.GetType()] = make(map[models.MetricName]interface{})
		m.storage[metric.GetType()][metric.GetName()] = int64(0)
	}

	switch string(metric.GetType()) {
	case models.MetricTypeCounter:
		oldVal := m.storage[metric.GetType()][metric.GetName()]
		m.storage[metric.GetType()][metric.GetName()] = oldVal.(int64) + metric.GetValue().(int64)

	case models.MetricTypeGauge:
		m.storage[metric.GetType()][metric.GetName()] = metric.GetValue()

	default:
		return models.ErrMetricTypeNotFound
	}

	return nil
}

func (m *MemStorage) GetMetricValue(metric models.GetMetrics) (interface{}, error) {
	m.RLock()
	defer m.RUnlock()

	if _, ok := m.storage[metric.Type]; !ok {
		return nil, models.ErrMetricTypeNotFound
	}

	val, ok := m.storage[metric.Type][metric.Name]
	if !ok {
		return nil, models.ErrMetricNameNotFound
	}

	return val, nil
}

func (m *MemStorage) GetAllMetrics() (map[models.MetricType]map[models.MetricName]interface{}, error) {
	if len(m.storage) == 0 {
		return nil, models.ErrStorageIsEmpty
	}
	return m.storage, nil
}
