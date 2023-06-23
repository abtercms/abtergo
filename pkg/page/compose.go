package page

import (
	"go.uber.org/zap"

	"github.com/abtergo/abtergo/libs/repo"
)

func CreateHandler(logger *zap.Logger) (*Handler, Repo) {
	r := repo.NewInMemoryRepo[Page]()
	u := NewUpdater()
	srv := NewService(r, u, logger)

	return NewHandler(srv, logger), r
}
