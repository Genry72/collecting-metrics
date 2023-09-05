package filestorage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/models"
)

func (fs *FileStorage) GetAllMetrics(ctx context.Context) ([]*models.Metric, error) {
	fs.mx.RLock()
	defer fs.mx.RUnlock()
	metrics := make([]*models.Metric, 0)
	for fs.reader.Scan() {
		if err := checkContext(ctx); err != nil {
			return nil, fmt.Errorf("GetAllMetrics: %w", err)
		}
		b := fs.reader.Bytes()
		metric := models.Metric{}
		if err := json.Unmarshal(b, &metric); err != nil {
			return nil, fmt.Errorf("GetAllMetrics.Unmarshal: %w", err)
		}
		metrics = append(metrics, &metric)
	}

	return metrics, nil
}
