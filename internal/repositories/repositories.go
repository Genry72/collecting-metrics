package repositories

import (
	"context"
	"github.com/Genry72/collecting-metrics/internal/models"
	"github.com/Genry72/collecting-metrics/internal/repositories/fileStorage"
)

type Repositories interface {
	SetMetric(ctx context.Context, metric *models.Metrics) (*models.Metrics, error)
	SetAllMetrics(ctx context.Context, metrics []models.Metrics) error
	GetMetricValue(ctx context.Context, metric *models.Metrics) (*models.Metrics, error)
	GetAllMetrics(ctx context.Context) ([]models.Metrics, error)
}

type PermanentStorage interface {
	SetAllMetrics(context.Context, []models.Metrics) error
	GetAllMetrics(ctx context.Context) ([]models.Metrics, error)
	Start() error
	Stop() error
	IsStarted() bool
	GetConfig() *fileStorage.StorageConf
}
