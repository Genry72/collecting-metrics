package usecases

import (
	"context"
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/models"
	"github.com/Genry72/collecting-metrics/internal/repositories"

	"strconv"
	"strings"
)

type ServerUc struct {
	memStorage repositories.Repositories
}

func NewServerUc(repo repositories.Repositories) *ServerUc {
	return &ServerUc{
		memStorage: repo,
	}
}

func (uc *ServerUc) SetMetric(ctx context.Context, metric *models.UpdateMetrics) error {

	switch metric.Type {
	case models.MetricTypeGauge:
		val, err := strconv.ParseFloat(metric.Value, 64)
		if err != nil {
			return fmt.Errorf("%w: %s", models.ErrParseValue, err.Error())
		}
		if err := uc.memStorage.SetMetricGauge(ctx, metric.Name, val); err != nil {
			return err
		}

	case models.MetricTypeCounter:
		val, err := strconv.ParseInt(metric.Value, 10, 64)
		if err != nil {
			return fmt.Errorf("%w: %s", models.ErrParseValue, err.Error())
		}
		if err := uc.memStorage.SetMetricCounter(ctx, metric.Name, val); err != nil {
			return err
		}

	default:
		return fmt.Errorf("%w: %s", models.ErrBadMetricType, metric.Type)
	}

	return nil
}

func (uc *ServerUc) GetMetricValue(ctx context.Context, metric models.GetMetrics) (interface{}, error) {
	switch metric.Type {

	case models.MetricTypeGauge:
		return uc.memStorage.GetMetricValueGauge(ctx, metric.Name)

	case models.MetricTypeCounter:
		return uc.memStorage.GetMetricValueCounter(ctx, metric.Name)

	default:
		return nil, fmt.Errorf("%w: %s", models.ErrBadMetricType, metric.Type)
	}

}

func (uc *ServerUc) GetAllMetrics(ctx context.Context) (string, error) {
	mapa, err := uc.memStorage.GetAllMetrics(ctx)
	if err != nil {
		return "", err
	}

	sb := strings.Builder{}

	for k, v := range mapa {
		sb.WriteString(fmt.Sprintf("%s : %v\n", k, v))
	}

	return sb.String(), nil
}
