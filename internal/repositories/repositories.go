package repositories

import (
	"context"
	"github.com/Genry72/collecting-metrics/internal/models"
	"github.com/Genry72/collecting-metrics/internal/repositories/filestorage"
)

type Repositories interface {
	// SetMetric добавление/изменение метрики
	SetMetric(ctx context.Context, metrics ...*models.Metric) error
	// SetAllMetrics добавление/изменение метрик
	SetAllMetrics(ctx context.Context, metrics []*models.Metric) error
	// GetMetricValue получение значения метрики
	GetMetricValue(ctx context.Context, metricType models.MetricType,
		metricName models.MetricName) (*models.Metric, error)
	// GetAllMetrics получение всех метрик
	GetAllMetrics(ctx context.Context) ([]*models.Metric, error)
}

type PermanentStorage interface {
	// SetAllMetrics добавление/изменение метрик
	SetAllMetrics(context.Context, []*models.Metric) error
	// GetAllMetrics получение всех метрик
	GetAllMetrics(ctx context.Context) ([]*models.Metric, error)
	// Stop Остановка работы с файловым хранилищем
	Stop()
	// GetConfig получение конфигурации
	GetConfig() *filestorage.StorageConf
}

type DatabaseStorage interface {
	// Ping Проверка доступности базы данных
	Ping() error
	// Stop Остановка работы с базой даных
	Stop()
}
