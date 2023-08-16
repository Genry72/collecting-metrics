package logger

import (
	"go.uber.org/zap"
	"log"
)

func NewZapLogger(level string) *zap.Logger {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		log.Fatal(err)
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return logger
}
