package usecases

import (
	"errors"
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/repositories"
)

const (
	metricTypeGauge   = "gauge"
	metricTypeCounter = "counter"
)

var (
	ErrBadMetricType = errors.New("bad metric type")
	ErrParseValue    = errors.New("fail parse metric value")
)

type ServerUc struct {
	memStorage repositories.Repositories
}

func NewServerUc(repo repositories.Repositories) *ServerUc {
	return &ServerUc{
		memStorage: repo,
	}
}

func (uc *ServerUc) SetMetric(typ, name, value string) error {
	var (
		m   repositories.Metric
		err error
	)

	switch typ {
	case metricTypeGauge:
		m, err = newGauge(typ, name, value)
		if err != nil {
			return err
		}

	case metricTypeCounter:
		m, err = newCounter(typ, name, value)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("%w: %s", ErrBadMetricType, typ)
	}

	if err := uc.memStorage.SetMetric(m); err != nil {
		return err
	}

	return nil
}
