package repositories

type Repositories interface {
	SetMetric(metric Metric) error
}

type Metric interface {
	GetType() string
	GetName() string
	GetValue() interface{}
}
