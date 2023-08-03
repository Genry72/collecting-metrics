package models

import "errors"

var (
	ErrStorageIsEmpty     = errors.New("metricStorage is empry")
	ErrMetricTypeNotFound = errors.New("metric type not found")
	ErrMetricNameNotFound = errors.New("metric name not found")
	ErrBadMetricType      = errors.New("bad metric type")
	ErrParseValue         = errors.New("fail parse metric value")
	ErrOnlyPost           = errors.New("only POST requests are allowed")
	ErrFormatURL          = errors.New("only /update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ> format are allowed")
)
