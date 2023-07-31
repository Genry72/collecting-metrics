package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	metricTypeGauge   = "gauge"
	metricTypeCounter = "counter"
)

type MyHandler struct {
	MemStorage
}

type MemStorage struct {
	storage map[MetricType]map[MetricName]MetricsValue
}

type MetricType string
type MetricName string
type MetricsValue interface{}

type Metric struct {
	MetricType
	MetricName
	MetricsValue
}

func main() {
	var h MyHandler
	ms := NewStorage(
		Metric{
			MetricType: "gauge",
		},
		Metric{
			MetricType: "counter",
		})

	h.storage = ms.storage

	mux := http.NewServeMux()
	mux.Handle(`/update/`, middlewareSerMetrics(http.HandlerFunc(h.setMetrics)))
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}

func (h MyHandler) setMetrics(w http.ResponseWriter, r *http.Request) {
	const (
		metricType = iota + 2
		metricName
		metricValue
	)

	urlSlice := strings.Split(r.URL.Path, "/")

	if err := h.setMetric(Metric{
		MetricType:   MetricType(urlSlice[metricType]),
		MetricName:   MetricName(urlSlice[metricName]),
		MetricsValue: urlSlice[metricValue],
	}); err != nil {
		http.Error(w, err.Error(),
			http.StatusBadRequest)
		return
	}

	if _, err := w.Write([]byte(fmt.Sprintf("%+v\n", h.storage))); err != nil {
		fmt.Println(err)
	}
}

func middlewareSerMetrics(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
			return
		}

		urlSlice := strings.Split(r.URL.Path, "/")
		if len(urlSlice) != 5 {
			http.Error(w, "Only /update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ> format are allowed!",
				http.StatusNotFound)
			return
		}

		//for _, p := range urlSlice {
		//	if p == "" {
		//		http.Error(w, "p!",
		//			http.StatusBadRequest)
		//		return
		//	}
		//}

		next.ServeHTTP(w, r)
	})
}

func NewStorage(metric ...Metric) *MemStorage {
	storage := make(map[MetricType]map[MetricName]MetricsValue)
	for _, m := range metric {
		if storage[m.MetricType] == nil {
			storage[m.MetricType] = make(map[MetricName]MetricsValue)
		}
	}
	return &MemStorage{
		storage: storage,
	}
}

func (m MemStorage) setMetric(metric Metric) error {
	fmt.Printf("%+v\n", metric)
	switch metric.MetricType {
	case metricTypeGauge:
		val, err := strconv.ParseFloat(metric.MetricsValue.(string), 64)
		if err != nil {
			return err
		}
		metric.MetricsValue = val
	case metricTypeCounter:
		val, err := strconv.ParseInt(metric.MetricsValue.(string), 10, 64)
		if err != nil {
			return err
		}
		metric.MetricsValue = val
	default:
		return fmt.Errorf("bad metruc type")
	}

	m.storage[metric.MetricType][metric.MetricName] = metric.MetricsValue
	fmt.Println(m.storage)
	return nil
}
