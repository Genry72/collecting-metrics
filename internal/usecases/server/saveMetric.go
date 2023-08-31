package server

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"time"
)

func (uc *Server) RunSaveToPermanentStorage(ctx context.Context) {
	fsconf := uc.permanentStorage.GetConfig()
	// Запускаем периодическое сохранение метрик в файл
	go func() {
		for {
			select {
			case <-ctx.Done():
				uc.log.Info("Stop RunSaveToPermanentStorage process")
				return
			default:
			}
			if err := uc.SaveToPermanentStorage(ctx); err != nil {
				uc.log.Error(err.Error())
				return
			}

			time.Sleep(time.Duration(fsconf.StoreInterval) * time.Second)
		}
	}()
}

// LoadMetricFromPermanentStore Загружаем метрики в memstorage
func (uc *Server) LoadMetricFromPermanentStore(ctx context.Context) error {

	metrics, err := uc.permanentStorage.GetAllMetrics(ctx)
	if err != nil {
		return fmt.Errorf("LoadMetricFromPermanentStore.GetAllMetrics: %w", err)
	}

	if err := uc.storage.SetAllMetrics(ctx, metrics); err != nil {
		return fmt.Errorf("LoadMetricFromPermanentStore.SetAllMetrics: %w", err)
	}

	uc.log.Info("metrics success loaded from start", zap.Int("count", len(metrics)))

	return nil
}

func (uc *Server) SaveToPermanentStorage(ctx context.Context) error {
	metrics, err := uc.storage.GetAllMetrics(ctx)
	if err != nil {
		return fmt.Errorf("saveToPermanentStorage.GetAllMetrics: %w", err)
	}

	if err := uc.permanentStorage.SetAllMetrics(ctx, metrics); err != nil {
		return fmt.Errorf("saveToPermanentStorage.SetAllMetrics: %w", err)
	}
	return nil
}
