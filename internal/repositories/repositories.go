package repositories

import (
	"context"
	"github.com/Genry72/collecting-metrics/internal/models"
)

type Repositories interface {
	SetMetricCounter(ctx context.Context, name models.MetricName, value int64) error
	SetMetricGauge(ctx context.Context, name models.MetricName, value float64) error
	GetMetricValueCounter(ctx context.Context, name models.MetricName) (int64, error)
	GetMetricValueGauge(ctx context.Context, name models.MetricName) (float64, error)
	GetAllMetrics(ctx context.Context) (map[models.MetricName]interface{}, error)
}

type Metric interface {
	GetType() models.MetricType
	GetName() models.MetricName
	GetValue() interface{}
}
