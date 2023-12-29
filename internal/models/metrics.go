package models

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
