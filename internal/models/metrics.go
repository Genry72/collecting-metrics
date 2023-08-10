package models

const (
	MetricTypeGauge   MetricType = "gauge"
	MetricTypeCounter MetricType = "counter"
)

type UpdateMetrics struct {
	Name  MetricName `uri:"name" binding:"required"`
	Type  MetricType `uri:"type" binding:"required"`
	Value string     `uri:"value" binding:"required"`
}

type GetMetrics struct {
	Name MetricName `uri:"name" binding:"required"`
	Type MetricType `uri:"type" binding:"required"`
}

type MetricType string

type MetricName string
