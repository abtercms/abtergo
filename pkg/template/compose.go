package template

import (
	"log/slog"

	"github.com/abtergo/abtergo/libs/repo"
)

func CreateHandler(logger *slog.Logger) *Handler {
	repository := repo.NewInMemoryRepo[Template]()
	srv := NewService(repository, logger)
	handler := NewHandler(srv, logger)

	return handler
}
