package main

import (
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/usecases"
	"time"
)

var (
	flagEndpointServer string // endpoint сервера
	flagReportInterval int    // частота оправки метрик в секундах
	flagPollInterval   int    // частота обновления метрик
)

const (
	envEndpoint       = "ADDRESS"
	envreportInterval = "REPORT_INTERVAL"
	envPollInterval   = "POLL_INTERVAL"
)

func main() {
	fmt.Println("start agent")

	metrics := usecases.NewMetrics()

	// обрабатываем аргументы командной строки
	parseFlags()
	// Запускаем обновление раз в 2 секунты
	metrics.Update(time.Duration(flagPollInterval) * time.Second)

	agent := usecases.NewAgent("http://" + flagEndpointServer)

	// Запускаем отправку метрик раз 10 секунд
	agent.SendMetrics(metrics, time.Duration(flagReportInterval)*time.Second)

}
