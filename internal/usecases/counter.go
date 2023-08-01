package usecases

import (
	"fmt"
	"strconv"
)

type counter struct {
	typ   string
	name  string
	value int64
}

func newCounter(typ, name, value string) (*counter, error) {
	val, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrParseValue, err.Error())
	}

	return &counter{
		typ:   typ,
		name:  name,
		value: val,
	}, nil
}

func (m *counter) GetType() string {
	return m.typ
}

func (m *counter) GetName() string {
	return m.name
}

func (m *counter) GetValue() interface{} {
	return m.value
}
