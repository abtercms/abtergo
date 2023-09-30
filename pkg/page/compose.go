package page

import (
	"log/slog"

	"github.com/abtergo/abtergo/libs/repo"
)

func CreateHandler(logger *slog.Logger) (*Handler, Repo) {
	r := repo.NewInMemoryRepo[Page]()
	u := NewUpdater()
	srv := NewService(r, u, logger)

	return NewHandler(srv, logger), r
}
