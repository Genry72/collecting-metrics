package main

import (
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/handlers"
	"github.com/Genry72/collecting-metrics/internal/logger"
	"github.com/Genry72/collecting-metrics/internal/repositories/memstorage"
	"github.com/Genry72/collecting-metrics/internal/usecases/server"
)

var flagRunAddr string

const envRunAddr = "ADDRESS"

func main() {
	zapLogger := logger.NewZapLogger("info")

	defer func() {
		if err := zapLogger.Sync(); err != nil {
			fmt.Println(err)
		}
	}()

	zapLogger.Info("start server")
	repo := memstorage.NewMemStorage()

	uc := server.NewServerUc(repo)

	h := handlers.NewServer(uc, zapLogger)

	// обрабатываем аргументы командной строки
	parseFlags()

	if err := h.RunServer(flagRunAddr); err != nil {
		zapLogger.Fatal(err.Error())
		return
	}
}
