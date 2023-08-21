package repositories

import (
	"context"
	"github.com/Genry72/collecting-metrics/internal/models"
)

type Repositories interface {
	SetMetric(ctx context.Context, metric *models.Metrics) (*models.Metrics, error)
	GetMetricValue(ctx context.Context, metric *models.Metrics) (*models.Metrics, error)
	GetAllMetrics(ctx context.Context) (map[models.MetricName]interface{}, error)
}
