package main

import (
	"flag"
	"os"
	"strconv"
)

func parseFlags() {
	// регистрируем переменную flagEndpointServer
	// как аргумент -a со значением :8080 по умолчанию
	flag.StringVar(&flagEndpointServer, "a", ":8080", "address and port to run server")
	flag.IntVar(&flagReportInterval, "r", 10, "report interval")
	flag.IntVar(&flagPollInterval, "p", 2, "poll interval")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	if endpoint := os.Getenv(envEndpoint); endpoint != "" {
		flagEndpointServer = endpoint
	}

	if reportInterval := os.Getenv(envreportInterval); reportInterval != "" {
		if ri, err := strconv.ParseInt(reportInterval, 10, 64); err == nil {
			flagReportInterval = int(ri)
		}
	}

	if pollInterval := os.Getenv(envPollInterval); pollInterval != "" {
		if pi, err := strconv.ParseInt(pollInterval, 10, 64); err == nil {
			flagPollInterval = int(pi)
		}
	}

}
