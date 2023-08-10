package main

import (
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/handlers"
	"github.com/Genry72/collecting-metrics/internal/repositories/memstorage"
	"github.com/Genry72/collecting-metrics/internal/usecases/server"
	"log"
)

var flagRunAddr string

const envRunAddr = "ADDRESS"

func main() {
	fmt.Println("start server")
	repo := memstorage.NewMemStorage()

	uc := server.NewServerUc(repo)

	h := handlers.NewServer(uc)

	// обрабатываем аргументы командной строки
	parseFlags()

	if err := h.RunServer(flagRunAddr); err != nil {
		log.Println(err)
		return
	}
}
