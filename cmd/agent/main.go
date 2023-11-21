package main

import (
	"context"
	"github.com/Genry72/collecting-metrics/internal/logger"
	"github.com/Genry72/collecting-metrics/internal/usecases/agent"
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
)

const (
	envEndpoint       = "ADDRESS"
	envreportInterval = "REPORT_INTERVAL"
	envPollInterval   = "POLL_INTERVAL"
	envKeyHash        = "KEY"
	envRateLimit      = "RATE_LIMIT"
)

func main() {
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

	agentUc := agent.NewAgent("http://"+flagEndpointServer, zapLogger, keyHash, flagRateLimit)

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
