package template

import (
	"go.uber.org/zap"

	"github.com/abtergo/abtergo/libs/repo"
)

func CreateHandler(logger *zap.Logger) *Handler {
	repository := repo.NewInMemoryRepo[Template]()
	srv := NewService(repository, logger)
	handler := NewHandler(srv, logger)

	return handler
}
