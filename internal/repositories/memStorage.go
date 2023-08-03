package repositories

import (
	"sync"
)

type MemStorage struct {
	sync.Mutex
	storage map[string]map[string]interface{}
}

func NewMemStorage() *MemStorage {
	storage := make(map[string]map[string]interface{})

	return &MemStorage{
		storage: storage,
	}
}

func (m *MemStorage) SetMetric(metric Metric) error {
	m.Lock()

	if m.storage[metric.GetType()] == nil {
		m.storage[metric.GetType()] = make(map[string]interface{})
	}

	m.storage[metric.GetType()][metric.GetName()] = metric.GetValue()

	m.Unlock()

	return nil
}
