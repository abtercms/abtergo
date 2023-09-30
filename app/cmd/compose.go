package main

import (
	"log/slog"

	"github.com/adelowo/onecache"
	"github.com/adelowo/onecache/memory"
)

func createLogger() *slog.Logger {
	return slog.Default()
}

func createCache() onecache.Store {
	return memory.New()
}
