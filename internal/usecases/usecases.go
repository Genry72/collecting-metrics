package usecases

import (
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/models"
	"github.com/Genry72/collecting-metrics/internal/repositories"
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

func (uc *ServerUc) SetMetric(metric *models.UpdateMetrics) error {
	var (
		m   repositories.Metric
		err error
	)

	switch metric.Type {
	case models.MetricTypeGauge:
		m, err = newGauge(metric.Type, metric.Name, metric.Value)
		if err != nil {
			return err
		}

	case models.MetricTypeCounter:
		m, err = newCounter(metric.Type, metric.Name, metric.Value)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("%w: %s", models.ErrBadMetricType, metric.Type)
	}

	if err := uc.memStorage.SetMetric(m); err != nil {
		return err
	}

	return nil
}

func (uc *ServerUc) GetMetricValue(metric models.GetMetrics) (interface{}, error) {
	return uc.memStorage.GetMetricValue(metric)
}

func (uc *ServerUc) GetAllMetrics() (string, error) {
	mapa, err := uc.memStorage.GetAllMetrics()
	if err != nil {
		return "", err
	}

	sb := strings.Builder{}

	for _, v := range mapa {
		for kk, vv := range v {
			sb.WriteString(fmt.Sprintf("%s : %v\n", kk, vv))
		}
	}
	return sb.String(), nil
}
