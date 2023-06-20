package main

import (
	"github.com/adelowo/onecache"
	"github.com/adelowo/onecache/memory"
	"go.uber.org/zap"
)

func createLogger() *zap.Logger {
	return zap.Must(zap.NewProduction())
}

func createCache() onecache.Store {
	return memory.New()
}
