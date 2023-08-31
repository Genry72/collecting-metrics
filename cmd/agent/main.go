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
)

const (
	envEndpoint       = "ADDRESS"
	envreportInterval = "REPORT_INTERVAL"
	envPollInterval   = "POLL_INTERVAL"
)

func main() {

	zapLogger := logger.NewZapLogger("info")

	defer func() {
		if err := zapLogger.Sync(); err != nil {
			fmt.Println(err)
		}
	}()

	zapLogger.Info("start agent")

	metrics := agent.NewMetrics()

	// обрабатываем аргументы командной строки
	parseFlags()
	// Запускаем обновление раз в 2 секунты
	metrics.Update(time.Duration(flagPollInterval) * time.Second)

	agentUc := agent.NewAgent("http://"+flagEndpointServer, zapLogger)

	// Запускаем отправку метрик раз 10 секунд
	agentUc.SendMetrics(metrics, time.Duration(flagReportInterval)*time.Second)

}
