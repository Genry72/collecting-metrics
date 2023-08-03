package usecases

import (
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/models"
	"strconv"
)

type gauge struct {
	typ   models.MetricType
	name  models.MetricName
	value float64
}

func newGauge(typ models.MetricType, name models.MetricName, value string) (*gauge, error) {
	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", models.ErrParseValue, err.Error())
	}

	return &gauge{
		typ:   typ,
		name:  name,
		value: val,
	}, nil
}

func (m *gauge) GetType() models.MetricType {
	return m.typ
}

func (m *gauge) GetName() models.MetricName {
	return m.name
}

func (m *gauge) GetValue() interface{} {
	return m.value
}
