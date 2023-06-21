package redirect

import (
	"go.uber.org/zap"

	"github.com/abtergo/abtergo/libs/repo"
)

func CreateHandler(logger *zap.Logger) *Handler {
	r := repo.NewInMemoryRepo[Redirect]()
	srv := NewService(r, logger)

	return NewHandler(srv, logger)
}
