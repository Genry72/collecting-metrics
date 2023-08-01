package usecases

import (
	"fmt"
	"strconv"
)

type gauge struct {
	typ   string
	name  string
	value float64
}

func newGauge(typ, name, value string) (*gauge, error) {
	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrParseValue, err.Error())
	}

	return &gauge{
		typ:   typ,
		name:  name,
		value: val,
	}, nil
}

func (m *gauge) GetType() string {
	return m.typ
}

func (m *gauge) GetName() string {
	return m.name
}

func (m *gauge) GetValue() interface{} {
	return m.value
}
