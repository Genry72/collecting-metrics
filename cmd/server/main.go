package main

import (
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/handlers"
	"github.com/Genry72/collecting-metrics/internal/repositories"
	"github.com/Genry72/collecting-metrics/internal/usecases"
	"log"
)

var flagRunAddr string

const envRunAddr = "ADDRESS"

func main() {
	fmt.Println("start server")
	repo := repositories.NewMemStorage()

	uc := usecases.NewServerUc(repo)

	h := handlers.NewServer(uc)

	// обрабатываем аргументы командной строки
	parseFlags()

	if err := h.RunServer(flagRunAddr); err != nil {
		log.Println(err)
		return
	}
}
