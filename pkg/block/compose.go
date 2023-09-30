package block

import (
	"log/slog"

	"github.com/abtergo/abtergo/libs/repo"
)

func CreateHandler(logger *slog.Logger) *Handler {
	r := repo.NewInMemoryRepo[Block]()
	srv := NewService(r, logger)

	return NewHandler(srv, logger)
}
