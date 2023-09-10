package models

import (
	"errors"
)

var (
	ErrStorageIsEmpty     = errors.New("metricStorage is empry")
	ErrMetricTypeNotFound = errors.New("metric type not found")
	ErrMetricNameNotFound = errors.New("metric name not found")
	ErrBadMetricType      = errors.New("bad metric type")
	ErrParseValue         = errors.New("fail parse metric value")
	ErrFormatURL          = errors.New("only /update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ> format are allowed")
	ErrBadBody            = errors.New("bad format body")
	ErrBadMetricValue     = errors.New("bad metric value")
	ErrDeadlineContext    = errors.New("deadline Context")
)

// RetryError тип ошибки, при которой требуется переповтор запросов
type RetryError struct {
	err error
}

func (re *RetryError) Error() string {
	return re.err.Error()
}

func NewRetryError(err error) error {
	return &RetryError{
		err: err,
	}
}
