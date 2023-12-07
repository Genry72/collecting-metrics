package agent

import (
	"fmt"
	"github.com/Genry72/collecting-metrics/helpers"
	"github.com/Genry72/collecting-metrics/internal/models"
	"github.com/fatih/structs"
	"testing"
)

func BenchmarkGetMetrics(b *testing.B) {
	m := &Metrics{
		gauge:   &gaugeRunTimeMetrics{Alloc: 3},
		counter: &counterMetrics{},
	}

	b.ResetTimer()
	b.Run("new", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = m.getMetricsNew()
		}
	})
	b.Run("old", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = m.getMetricsOld()
		}
	})

}

func (m *Metrics) getMetricsNew() ([]*models.Metric, error) {
	m.sm.RLock()
	defer m.sm.RUnlock()
	gaugeMetricData := helpers.StructToMap(m.gauge)
	counterMetricsData := helpers.StructToMap(m.counter)

	result := make([]*models.Metric, len(gaugeMetricData)+len(counterMetricsData))
	i := 0

	for metricName, value := range gaugeMetricData {
		v, err := fromInterfaceGauge(value)
		if err != nil {
			return nil, fmt.Errorf("fromInterfaceGauge: %w", err)
		}
		result[i] = &models.Metric{
			ID:        models.MetricName(metricName),
			MType:     models.MetricTypeGauge,
			Delta:     nil,
			Value:     v,
			ValueText: fmt.Sprint(value),
		}
		i++
	}

	for metricName, value := range counterMetricsData {
		v, ok := value.(int64)
		if !ok {
			return nil, fmt.Errorf("value.(int64): %w", models.ErrBadMetricValue)
		}

		result[i] = &models.Metric{
			ID:        models.MetricName(metricName),
			MType:     models.MetricTypeCounter,
			Delta:     &v,
			Value:     nil,
			ValueText: fmt.Sprint(value),
		}
		i++
	}

	return result, nil
}

// Получение метрик
func (m *Metrics) getMetricsOld() ([]*models.Metric, error) {
	m.sm.RLock()
	defer m.sm.RUnlock()
	gaugeMetricData := structs.Map(m.gauge)
	counterMetricsData := structs.Map(m.counter)

	result := make([]*models.Metric, 0, len(gaugeMetricData)+len(counterMetricsData))

	for metricName, value := range gaugeMetricData {
		v, err := fromInterfaceGauge(value)
		if err != nil {
			return nil, fmt.Errorf("fromInterfaceGauge: %w", err)
		}
		result = append(result, &models.Metric{
			ID:        models.MetricName(metricName),
			MType:     models.MetricTypeGauge,
			Delta:     nil,
			Value:     v,
			ValueText: fmt.Sprint(value),
		})
	}

	for metricName, value := range counterMetricsData {
		v, ok := value.(int64)
		if !ok {
			return nil, fmt.Errorf("value.(int64): %w", models.ErrBadMetricValue)
		}

		result = append(result, &models.Metric{
			ID:        models.MetricName(metricName),
			MType:     models.MetricTypeCounter,
			Delta:     &v,
			Value:     nil,
			ValueText: fmt.Sprint(value),
		})
	}

	return result, nil
}
