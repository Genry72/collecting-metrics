package filestorage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/models"
	"go.uber.org/zap"
)

func (fs *FileStorage) SetAllMetrics(ctx context.Context, metrics []*models.Metric) error {
	fs.mx.Lock()
	if err := fs.file.Truncate(0); err != nil {
		fs.mx.Unlock()
		return fmt.Errorf("fs.file.Truncate: %w", err)
	}

	fs.mx.Unlock()

	for _, metric := range metrics {
		if err := checkContext(ctx); err != nil {
			return fmt.Errorf("checkContext: %w", err)
		}
		if err := fs.write(ctx, metric); err != nil {
			return fmt.Errorf("fs.write: %w", err)
		}
	}

	fs.log.Info("Write metrics success", zap.Int("count", len(metrics)))

	return nil
}

func (fs *FileStorage) write(ctx context.Context, metric *models.Metric) error {
	fs.mx.Lock()
	defer fs.mx.Unlock()

	if err := checkContext(ctx); err != nil {
		return fmt.Errorf("checkContext: %w", err)
	}

	data, err := json.Marshal(metric)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	// записываем событие в буфер
	if _, err := fs.writer.Write(data); err != nil {
		return fmt.Errorf("fs.writer.Write: %w", err)
	}

	// добавляем перенос строки
	if err := fs.writer.WriteByte('\n'); err != nil {
		return err
	}

	// записываем буфер в файл
	return fs.writer.Flush()
}
