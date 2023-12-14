package main

import (
	"context"
	"fmt"
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

var (
	flagRunAddr         string
	flagStoreInterval   int
	flagFileStoragePath string
	flagRestore         bool
	flagPgDsn           string
	flagKeyHash         string
	flagCryptKey        string // Путь до файла с приватным ключом
)

// Информация о сборке
var (
	buildVersion string
	buildDate    string
	buildCommit  string
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
	// Строка с адресом подключения к БД
	envPgDSN = "DATABASE_DSN"
	// Ключ для шифрования
	envKeyHash = "KEY"
	// Путь до файла с приватным ключом
	envCryptKey = "CRYPTO_KEY"
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
	parseFlags()

	zapLogger.Info("start server")

	ctxMain, mainStop := context.WithCancel(context.Background())

	defer mainStop()

	databaseStorage, err := postgre.NewPGStorage(flagPgDsn, zapLogger)
	if err != nil {
		zapLogger.Error("connect databaseStorage", zap.Error(err))
	} else {
		zapLogger.Info("connect to db success")
		defer databaseStorage.Stop()
	}

	var uc *server.Server

	if databaseStorage == nil {
		memStorage := memstorage.NewMemStorage(zapLogger)

		permStorConf := filestorage.NewPermanentStorageConf(flagStoreInterval, flagFileStoragePath, flagRestore)

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

	var keyHash *string

	if flagKeyHash != "" {
		keyHash = &flagKeyHash
	}

	go func() {
		if err := h.RunServer(flagRunAddr, keyHash, flagCryptKey, organization); err != nil {
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
