package models

const (
	MetricTypeGauge   MetricType = "gauge"
	MetricTypeCounter MetricType = "counter"
)

type MetricType string

type MetricName string

// Metrics структура body запроса/ ответа
type Metrics struct {
	ID        MetricName `json:"id" uri:"name" binding:"required" binding:"required"` // Имя метрики
	MType     MetricType `json:"type" uri:"type" binding:"required"`                  // Параметр, принимающий значение gauge или counter
	Delta     *int64     `json:"delta,omitempty"`                                     // Значение метрики в случае передачи counter
	Value     *float64   `json:"value,omitempty"`                                     // Значение метрики в случае передачи gauge
	ValueText string     `json:"-" uri:"value"`                                       // Значение метрики в случае передачи GET запросом
}
