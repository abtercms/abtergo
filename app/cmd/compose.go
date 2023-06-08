package main

import (
	"go.uber.org/zap"
)

func createLogger() *zap.Logger {
	return zap.Must(zap.NewProduction())
}
