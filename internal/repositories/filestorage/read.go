package filestorage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/models"
)

func (fs *FileStorage) GetAllMetrics(ctx context.Context) ([]models.Metrics, error) {
	if err := fs.Start(); err != nil {
		return nil, err
	}

	fs.mx.RLock()
	defer fs.mx.RUnlock()
	metrics := make([]models.Metrics, 0)
	for fs.reader.Scan() {
		b := fs.reader.Bytes()
		metric := models.Metrics{}
		if err := json.Unmarshal(b, &metric); err != nil {
			return nil, fmt.Errorf("GetAllMetrics.Unmarshal: %w", err)
		}
		metrics = append(metrics, metric)
	}

	return metrics, nil
}
