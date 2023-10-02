package main

import (
	"flag"
	"log"
	"os"
	"strconv"
)

func parseFlags() {
	// регистрируем переменную flagRunAddr
	// как аргумент -a со значением :8080 по умолчанию
	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")

	flag.IntVar(&flagStoreInterval, "i", 300, "интервал времени в секундах, по истечении которого"+
		" текущие показания сервера сохраняются на диск (по умолчанию 300 секунд, значение 0 делает запись синхронной)")

	flag.StringVar(&flagFileStoragePath, "f", "/tmp/metrics-db.json", "полное имя файла, куда сохраняются"+
		" текущие значения (по умолчанию /tmp/metrics-db.json, пустое значение отключает функцию записи на диск)")

	flag.BoolVar(&flagRestore, "r", true, "булево значение (true/false), определяющее, загружать или"+
		" нет ранее сохранённые значения из указанного файла при старте сервера (по умолчанию true)")

	flag.StringVar(&flagPgDsn, "d", "", "булево значение (true/false), определяющее, загружать или"+
		" нет ранее сохранённые значения из указанного файла при старте сервера (по умолчанию true)")
	flag.StringVar(&flagKeyHash, "k", "", "key for hash")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	if runAddr := os.Getenv(envRunAddr); runAddr != "" {
		flagRunAddr = runAddr
	}

	if value := os.Getenv(envStoreInterval); value != "" {
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		flagStoreInterval = int(v)
	}

	if value := os.Getenv(envFileStoragePath); value != "" {
		flagFileStoragePath = value
	}

	if value := os.Getenv(envRestore); value != "" {
		v, err := strconv.ParseBool(value)
		if err != nil {
			log.Fatal(err)
		}
		flagRestore = v
	}

	if value := os.Getenv(envPgDSN); value != "" {
		flagPgDsn = value
	}

	if key := os.Getenv(envKeyHash); key != "" {
		flagKeyHash = key
	}
}
