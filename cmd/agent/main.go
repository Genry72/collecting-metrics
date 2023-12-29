package main

import (
	"context"
	"fmt"
	"github.com/Genry72/collecting-metrics/cmd/agent/flags"
	"github.com/Genry72/collecting-metrics/internal/logger"
	"github.com/Genry72/collecting-metrics/internal/usecases/agent"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

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

	zapLogger.Info("start agent")

	// обрабатываем аргументы командной строки
	conf, err := flags.ParseFlag()
	if err != nil {
		zapLogger.Fatal("flags.ParseFlag", zap.Error(err))
	}

	metrics := agent.NewMetrics()

	ctx, cansel := context.WithCancel(context.Background())

	// Запускаем обновление раз в 2 секунты
	go func() {
		metrics.Update(ctx, time.Duration(*conf.PollInterval)*time.Second)
	}()

	if conf.Address == nil || *conf.Address == "" {
		zapLogger.Fatal("empty endpoint")
	}

	agentUc, err := agent.NewAgent("http://"+*conf.Address,
		conf.GrpcAddress, zapLogger, conf.KeyHash, conf.CryptoKey, conf.RateLimit)

	if err != nil {
		zapLogger.Fatal("agent.NewAgent", zap.Error(err))
	}

	// Запускаем отправку метрик раз 10 секунд
	go func() {
		agentUc.SendMetrics(ctx, metrics, time.Duration(*conf.ReportInterval)*time.Second)
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
