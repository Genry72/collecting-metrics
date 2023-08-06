package main

import (
	"flag"
	"os"
)

func parseFlags() {
	// регистрируем переменную flagRunAddr
	// как аргумент -a со значением :8080 по умолчанию
	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	if runAddr := os.Getenv(envRunAddr); runAddr != "" {
		flagRunAddr = runAddr
	}
}
