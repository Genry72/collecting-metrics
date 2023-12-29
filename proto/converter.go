package proto

import (
	"github.com/Genry72/collecting-metrics/internal/models"
)

func MetricsDpToMetrics(metrics *Metrics) []*models.Metric {
	result := make([]*models.Metric, len(metrics.Metrics))

	for k := range metrics.Metrics {
		result[k] = &models.Metric{
			ID:        models.MetricName(metrics.Metrics[k].Id),
			MType:     models.MetricType(metrics.Metrics[k].Type),
			Delta:     &metrics.Metrics[k].Delta,
			Value:     &metrics.Metrics[k].Value,
			ValueText: metrics.Metrics[k].ValueText,
		}
	}

	return result
}

func MetricsToMetricsDp(metrics []*models.Metric) *Metrics {
	result := &Metrics{
		Metrics: make([]*Metric, len(metrics)),
	}

	for k := range metrics {
		result.Metrics[k] = &Metric{
			Id:        string(metrics[k].ID),
			Type:      string(metrics[k].MType),
			ValueText: metrics[k].ValueText,
		}

		if metrics[k].Delta != nil {
			result.Metrics[k].Delta = *metrics[k].Delta
		}

		if metrics[k].Value != nil {
			result.Metrics[k].Value = *metrics[k].Value
		}
	}

	return result
}
