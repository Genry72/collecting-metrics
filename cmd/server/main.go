package main

import (
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/handlers"
	"github.com/Genry72/collecting-metrics/internal/repositories"
	"github.com/Genry72/collecting-metrics/internal/usecases"
	"log"
)

func main() {
	fmt.Println("start server")
	repo := repositories.NewMemStorage()

	uc := usecases.NewServerUc(repo)

	h := handlers.NewServer(uc)

	if err := h.RunServer("8080"); err != nil {
		log.Println(err)
		return
	}

}
