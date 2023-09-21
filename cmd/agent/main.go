package main

import (
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/logger"
	"github.com/Genry72/collecting-metrics/internal/usecases/agent"
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
		if err := zapLogger.Sync(); err != nil {
			fmt.Println(err)
		}
	}()

	zapLogger.Info("start agent")

	metrics := agent.NewMetrics()

	// Запускаем обновление раз в 2 секунты
	metrics.Update(time.Duration(flagPollInterval) * time.Second)

	var keyHash *string

	if flagKeyHash != "" {
		keyHash = &flagKeyHash
	}

	agentUc := agent.NewAgent("http://"+flagEndpointServer, zapLogger, keyHash, flagRateLimit)

	// Запускаем отправку метрик раз 10 секунд
	agentUc.SendMetrics(metrics, time.Duration(flagReportInterval)*time.Second)

}
