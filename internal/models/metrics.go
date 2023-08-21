package models

const (
	MetricTypeGauge   MetricType = "gauge"
	MetricTypeCounter MetricType = "counter"
)

type MetricType string

type MetricName string

type SetGetMetrics struct {
	Name         MetricName
	Type         MetricType
	ValueGauge   *float64
	ValueCounter *int64
}

// text/plain

// SetMetricsText структура параметров запроса по добавлению значения метрики
type SetMetricsText struct {
	Name  MetricName `uri:"name" binding:"required"`
	Type  MetricType `uri:"type" binding:"required"`
	Value string     `uri:"value" binding:"required"`
}

// GetMetrics структура параметров запроса по получению значения метрики
type GetMetrics struct {
	Name MetricName `uri:"name" binding:"required"`
	Type MetricType `uri:"type" binding:"required"`
}

// JSON

// SetMetricsJson структура body запроса/ ответа в формате json
type Metrics struct {
	ID    string   `json:"id" binding:"required"`   // Имя метрики
	MType string   `json:"type" binding:"required"` // Параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"`         // Значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"`         // Значение метрики в случае передачи gauge
}
