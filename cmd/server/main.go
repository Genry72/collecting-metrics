package main

import (
	"github.com/Genry72/collecting-metrics/internal/handlers"
	"github.com/Genry72/collecting-metrics/internal/repositories"
	"github.com/Genry72/collecting-metrics/internal/usecases"
	"log"
)

func main() {
	repo := repositories.NewMemStorage()

	uc := usecases.New(repo)

	h := handlers.New(uc)

	if err := h.RunServer("8080"); err != nil {
		log.Fatal(err)
	}

}
