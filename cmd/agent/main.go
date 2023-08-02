package main

import (
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/usecases"
	"time"
)

func main() {
	fmt.Println("start agent")

	metrics := usecases.NewMetrics()

	// Запускаем обновление раз в 2 секунты
	metrics.Update(2 * time.Second)

	agent := usecases.NewAgent("http://localhost:8080")

	// Запускаем отправку метрик раз 10 секунд
	agent.SendMetrics(metrics, 10*time.Second)

}
