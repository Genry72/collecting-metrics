package models

import (
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/usecases/cryptor"
	jsoniter "github.com/json-iterator/go"
)

const (
	MetricTypeGauge   MetricType = "gauge"
	MetricTypeCounter MetricType = "counter"
)

type MetricType string

type MetricName string

// Metric структура body запроса/ ответа
type Metric struct {
	ID        MetricName `json:"id" uri:"name" binding:"required" db:"name"`   // Имя метрики
	MType     MetricType `json:"type" uri:"type" binding:"required" db:"type"` // Параметр, принимающий значение gauge или counter
	Delta     *int64     `json:"delta,omitempty" db:"delta"`                   // Значение метрики в случае передачи counter
	Value     *float64   `json:"value,omitempty" db:"value"`                   // Значение метрики в случае передачи gauge
	ValueText string     `json:"-" uri:"value"`                                // Значение метрики в случае передачи GET запросом
}

// Metrics список метрик
type Metrics []*Metric

// Encode Хеш SHA256 на основе ключа
func (m *Metrics) Encode(password string) (string, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	metricByte, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("json.Marshal: %w", err)
	}

	return cryptor.Encrypt(metricByte, password)
}
