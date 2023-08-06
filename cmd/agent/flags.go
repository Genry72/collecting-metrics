package main

import "flag"

func parseFlags() {
	// регистрируем переменную flagEndpointServer
	// как аргумент -a со значением :8080 по умолчанию
	flag.StringVar(&flagEndpointServer, "a", ":8080", "address and port to run server")
	flag.IntVar(&flagReportInterval, "r", 10, "report interval")
	flag.IntVar(&flagPollInterval, "p", 2, "poll interval")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
}
