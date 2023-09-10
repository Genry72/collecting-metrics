package repositories

import (
	"context"
	"github.com/Genry72/collecting-metrics/internal/models"
	"github.com/Genry72/collecting-metrics/internal/repositories/filestorage"
)

type Repositories interface {
	SetMetric(ctx context.Context, metrics ...*models.Metric) ([]*models.Metric, error)
	SetAllMetrics(ctx context.Context, metrics []*models.Metric) error
	GetMetricValue(ctx context.Context, metricType models.MetricType, metricName models.MetricName) (*models.Metric, error)
	GetAllMetrics(ctx context.Context) ([]*models.Metric, error)
}

type PermanentStorage interface {
	SetAllMetrics(context.Context, []*models.Metric) error
	GetAllMetrics(ctx context.Context) ([]*models.Metric, error)
	Stop()
	GetConfig() *filestorage.StorageConf
}

type DatabaseStorage interface {
	Ping() error
	Stop()
}
