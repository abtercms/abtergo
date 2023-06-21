package website

import (
	"github.com/adelowo/onecache"
	"go.uber.org/zap"

	"github.com/abtergo/abtergo/libs/templ"
)

func CreateHandler(logger *zap.Logger, cache onecache.Store) *Handler {
	cr := NewContentRetriever()
	tr := NewTemplateRetriever()
	p := templ.NewParser("block")
	r := templ.NewRenderer(p, nil, cache)
	srv := NewService(cr, tr, r)

	return NewHandler(srv, logger)
}
