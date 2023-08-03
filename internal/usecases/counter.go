package usecases

import (
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/models"
	"strconv"
)

type counter struct {
	typ   models.MetricType
	name  models.MetricName
	value int64
}

func newCounter(typ models.MetricType, name models.MetricName, value string) (*counter, error) {
	val, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: %s name: %s", models.ErrParseValue, err.Error(), name)
	}

	return &counter{
		typ:   typ,
		name:  name,
		value: val,
	}, nil
}

func (m *counter) GetType() models.MetricType {
	return m.typ
}

func (m *counter) GetName() models.MetricName {
	return m.name
}

func (m *counter) GetValue() interface{} {
	return m.value
}
