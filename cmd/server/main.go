package main

import (
	"context"
	"github.com/Genry72/collecting-metrics/internal/handlers"
	"github.com/Genry72/collecting-metrics/internal/logger"
	"github.com/Genry72/collecting-metrics/internal/repositories/fileStorage"
	"github.com/Genry72/collecting-metrics/internal/repositories/memstorage"
	"github.com/Genry72/collecting-metrics/internal/usecases/server"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	flagRunAddr         string
	flagStoreInterval   int
	flagFileStoragePath string
	flagRestore         bool
)

const (
	envRunAddr = "ADDRESS"
	/*
		Интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск
		(по умолчанию 300 секунд, значение 0 делает запись синхронной).
	*/
	envStoreInterval = "STORE_INTERVAL"
	/*
		Полное имя файла, куда сохраняются текущие значения (по умолчанию /tmp/metrics-db.json,
		пустое значение отключает функцию записи на диск)
	*/
	envFileStoragePath = "FILE_STORAGE_PATH"
	/*
		Булево значение (true/false), определяющее, загружать или нет ранее сохранённые значения
		из указанного файла при старте сервера (по умолчанию true)
	*/
	envRestore = "RESTORE"
)

func main() {
	zapLogger := logger.NewZapLogger("info")

	defer func() {
		_ = zapLogger.Sync()
	}()

	// обрабатываем аргументы командной строки
	parseFlags()

	zapLogger.Info("start server")

	repo := memstorage.NewMemStorage(zapLogger)

	permStorConf := fileStorage.NewPermanentStorageConf(flagStoreInterval, flagFileStoragePath, flagRestore)

	permStorage := fileStorage.NewFileStorage(permStorConf, zapLogger)

	uc := server.NewServerUc(repo, permStorage, zapLogger)

	// Загрузка метрик из файла при старте
	if permStorConf.Restore {
		if err := uc.LoadMetricFromPermanentStore(context.Background()); err != nil {
			zapLogger.Fatal(err.Error())
		}
	}

	h := handlers.NewServer(uc, zapLogger)

	go func() {
		if err := h.RunServer(flagRunAddr); err != nil {
			zapLogger.Fatal(err.Error())
		}
	}()

	// 	Запуск периодической отправки метрик в файл
	uc.RunSaveToPermanentStorage(context.Background())

	// Graceful shutdown block
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-quit

	zapLogger.Info("Graceful shutdown")

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)

	defer shutdown()

	if err := uc.SaveToPermanentStorage(ctx); err != nil {
		zapLogger.Error(err.Error())
	}
	zapLogger.Info("Success graceful shutdown")
}
