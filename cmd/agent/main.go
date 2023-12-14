package main

import (
	"context"
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/logger"
	"github.com/Genry72/collecting-metrics/internal/usecases/agent"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	flagEndpointServer string // endpoint сервера
	flagReportInterval int    // Частота оправки метрик в секундах
	flagPollInterval   int    // Частота обновления метрик
	flagKeyHash        string // Ключ для расчета HashSHA256
	flagRateLimit      uint64 // Количество одновременно исходящих запросов на сервер
	flagCryptKey       string // Путь до файла с приватным ключом
)

// Информация о сборке
var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

const (
	envEndpoint       = "ADDRESS"
	envreportInterval = "REPORT_INTERVAL"
	envPollInterval   = "POLL_INTERVAL"
	envKeyHash        = "KEY"
	envRateLimit      = "RATE_LIMIT"
	envCryptKey       = "CRYPTO_KEY"
)

func main() {
	// Печать информации о сборке
	fmt.Println("Build version:", printBuildInfo(buildVersion))
	fmt.Println("Build date:", printBuildInfo(buildDate))
	fmt.Println("Build commit:", printBuildInfo(buildCommit))

	// обрабатываем аргументы командной строки
	parseFlags()

	zapLogger := logger.NewZapLogger("info")

	defer func() {
		_ = zapLogger.Sync()
	}()

	zapLogger.Info("start agent")

	metrics := agent.NewMetrics()

	ctx, cansel := context.WithCancel(context.Background())

	// Запускаем обновление раз в 2 секунты
	go func() {
		metrics.Update(ctx, time.Duration(flagPollInterval)*time.Second)
	}()

	var keyHash *string

	if flagKeyHash != "" {
		keyHash = &flagKeyHash
	}

	agentUc, err := agent.NewAgent("http://"+flagEndpointServer, zapLogger, keyHash, flagCryptKey, flagRateLimit)
	if err != nil {
		zapLogger.Fatal("agent.NewAgent", zap.Error(err))
	}

	// Запускаем отправку метрик раз 10 секунд
	go func() {
		agentUc.SendMetrics(ctx, metrics, time.Duration(flagReportInterval)*time.Second)
	}()

	// Graceful shutdown block
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-quit

	cansel()
	// таймаут на завершение всех задач
	time.Sleep(time.Second)
}

func printBuildInfo(info string) string {
	if info != "" {
		return info
	}
	return "N/A"
}
