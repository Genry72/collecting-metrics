package main

import (
	"context"
	"fmt"
	"github.com/Genry72/collecting-metrics/cmd/server/flags"
	"github.com/Genry72/collecting-metrics/internal/handlers"
	"github.com/Genry72/collecting-metrics/internal/logger"
	"github.com/Genry72/collecting-metrics/internal/repositories/filestorage"
	"github.com/Genry72/collecting-metrics/internal/repositories/memstorage"
	"github.com/Genry72/collecting-metrics/internal/repositories/postgre"
	"github.com/Genry72/collecting-metrics/internal/usecases/server"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Организация, для выпоска сертификата и ключей
const organization = "vsemenovOrg"

// Информация о сборке
var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	// Печать информации о сборке
	fmt.Println("Build version:", printBuildInfo(buildVersion))
	fmt.Println("Build date:", printBuildInfo(buildDate))
	fmt.Println("Build commit:", printBuildInfo(buildCommit))

	zapLogger := logger.NewZapLogger("info")

	defer func() {
		_ = zapLogger.Sync()
	}()

	// обрабатываем аргументы командной строки
	conf, err := flags.ParseFlag()
	if err != nil {
		zapLogger.Fatal(err.Error())
	}

	zapLogger.Info("Started with config", zap.Any("config", conf))

	zapLogger.Info("start server")

	ctxMain, mainStop := context.WithCancel(context.Background())

	defer mainStop()

	databaseStorage, err := postgre.NewPGStorage(conf.DatabaseDsn, zapLogger)
	if err != nil {
		zapLogger.Error("connect databaseStorage", zap.Error(err))
	} else {
		zapLogger.Info("connect to db success")
		defer databaseStorage.Stop()
	}

	var uc *server.Server

	if databaseStorage == nil {
		memStorage := memstorage.NewMemStorage(zapLogger)

		if conf.StoreInterval == nil || conf.StoreFile == nil || conf.Restore == nil {
			zapLogger.Fatal("err load configurations")
		}

		permStorConf := filestorage.NewPermanentStorageConf(*conf.StoreInterval, *conf.StoreFile, *conf.Restore)

		fileStorage, err := filestorage.NewFileStorage(permStorConf, zapLogger)
		if err != nil {
			zapLogger.Error("start fileStorage", zap.Error(err))
		} else {
			zapLogger.Info("file storage success started")
			defer fileStorage.Stop()
		}

		uc = server.NewServerUc(memStorage, fileStorage, databaseStorage, zapLogger)

		// Загрузка метрик из файла при старте
		if permStorConf.Restore {
			if err := uc.LoadMetricFromPermanentStore(ctxMain); err != nil {
				zapLogger.Error(err.Error())
				return
			}
		}

	} else {
		uc = server.NewServerUc(databaseStorage, databaseStorage, databaseStorage, zapLogger)
	}

	h := handlers.NewServer(uc, zapLogger)

	go func() {
		if err := h.RunServer(conf.Address, conf.KeyHash, conf.CryptoKey, organization, conf.TrustedSubnet); err != nil {
			zapLogger.Fatal(err.Error())
		}
	}()

	// 	Запуск периодической отправки метрик в файл
	uc.RunSaveToPermanentStorage(ctxMain)

	// Graceful shutdown block
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-quit

	zapLogger.Info("Graceful shutdown")

	ctx, shutdown := context.WithTimeout(ctxMain, 5*time.Second)

	defer func() {
		mainStop()

		shutdown()

		zapLogger.Info("Success graceful shutdown")
	}()

	if err := uc.SaveToPermanentStorage(ctx); err != nil {
		zapLogger.Error(err.Error())
		return
	}
}

func printBuildInfo(info string) string {
	if info != "" {
		return info
	}
	return "N/A"
}
