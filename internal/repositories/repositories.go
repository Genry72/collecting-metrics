package repositories

import "github.com/Genry72/collecting-metrics/internal/models"

type Repositories interface {
	SetMetric(metric Metric) error
	GetMetricValue(metric models.GetMetrics) (interface{}, error)
	GetAllMetrics() (map[models.MetricType]map[models.MetricName]interface{}, error)
}

type Metric interface {
	GetType() models.MetricType
	GetName() models.MetricName
	GetValue() interface{}
}
